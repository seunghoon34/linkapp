package service

import (
	"errors"

	"github.com/seunghoon34/linkapp/backend/internal/model"
	"github.com/seunghoon34/linkapp/backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo     *repository.UserRepository
	linkRepo     *repository.LinkRepository
	chatroomRepo *repository.ChatroomRepository
}

func NewUserService(userRepo *repository.UserRepository, linkRepo *repository.LinkRepository, chatroomRepo *repository.ChatroomRepository) *UserService {
	return &UserService{
		userRepo:     userRepo,
		linkRepo:     linkRepo,
		chatroomRepo: chatroomRepo,
	}
}

func (s *UserService) CreateUser(user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.userRepo.Create(user)
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.userRepo.GetByID(id)

}

func (s *UserService) UpdateProfile(id string, profile model.Profile) error {
	return s.userRepo.UpdateProfile(id, profile)
}

func (s *UserService) UpdatePreferences(id string, preferences model.Preferences) error {
	return s.userRepo.UpdatePreferences(id, preferences)
}

func (s *UserService) UpdateUser(id string, username, email string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	user.Username = username
	user.Email = email

	return s.userRepo.Update(user)
}

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

func (s *UserService) AuthenticateUser(email, password string) (*model.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// Use the same error message for non-existent user to prevent email enumeration
			return nil, ErrInvalidCredentials
		}
		return nil, err // Return the original error for other database issues
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Use a generic error message to prevent timing attacks
		return nil, ErrInvalidCredentials
	}

	// Create a new user object without the password field
	authenticatedUser := &model.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		// Add other fields as necessary, but omit the password
	}

	return authenticatedUser, nil
}

func (s *UserService) SearchMatches(userID string, limit int) ([]*model.User, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return s.userRepo.SearchMatches(user, limit)
}

func (s *UserService) UpdateLocation(userID string, latitude, longitude float64) error {
	// Basic validation
	if latitude < -90 || latitude > 90 {
		return errors.New("invalid latitude")
	}
	if longitude < -180 || longitude > 180 {
		return errors.New("invalid longitude")
	}

	return s.userRepo.UpdateLocation(userID, latitude, longitude)
}

func (s *UserService) StartSearching(userID primitive.ObjectID) error {
	return s.userRepo.SetSearchingStatus(userID, true)
}

func (s *UserService) StopSearching(userID primitive.ObjectID) error {
	return s.userRepo.SetSearchingStatus(userID, false)
}

func (s *UserService) FindMatch(userID primitive.ObjectID) (*model.Link, error) {
	user, err := s.userRepo.GetByID(userID.Hex())
	if err != nil {
		return nil, err
	}

	if !user.IsSearching {
		return nil, errors.New("user is not in searching mode")
	}

	potentialMatch, err := s.userRepo.FindPotentialMatch(user)
	if err != nil {
		return nil, err
	}

	if potentialMatch == nil {
		return nil, errors.New("no potential match found")
	}

	link, err := s.linkRepo.CreateLink(user.ID, potentialMatch.ID)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.SetCurrentLink(user.ID, link.ID)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.SetCurrentLink(potentialMatch.ID, link.ID)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (s *UserService) RespondToLink(userID primitive.ObjectID, linkID primitive.ObjectID, accept bool) error {
	link, err := s.linkRepo.GetLink(linkID)
	if err != nil {
		return err
	}

	if link.UserAID != userID && link.UserBID != userID {
		return errors.New("user is not part of this link")
	}

	if accept {
		err = s.linkRepo.UpdateLinkStatus(linkID, model.LinkStatusAccepted)
		if err == nil {
			// Create a new chatroom
			_, err = s.chatroomRepo.CreateChatroom(linkID, link.UserAID, link.UserBID)
		}
	} else {
		err = s.linkRepo.UpdateLinkStatus(linkID, model.LinkStatusRejected)
		if err == nil {
			// If rejected, set both users back to searching
			s.userRepo.SetSearchingStatus(link.UserAID, true)
			s.userRepo.SetSearchingStatus(link.UserBID, true)
			s.userRepo.SetCurrentLink(link.UserAID, primitive.NilObjectID)
			s.userRepo.SetCurrentLink(link.UserBID, primitive.NilObjectID)
		}
	}

	return err
}

func (s *UserService) ExpireLinks() error {
	err := s.linkRepo.ExpireLinks()
	if err != nil {
		return err
	}

	expiredLinks, err := s.linkRepo.GetExpiredLinks()
	if err != nil {
		return err
	}

	for _, link := range expiredLinks {
		s.userRepo.SetSearchingStatus(link.UserAID, true)
		s.userRepo.SetSearchingStatus(link.UserBID, true)
		s.userRepo.SetCurrentLink(link.UserAID, primitive.NilObjectID)
		s.userRepo.SetCurrentLink(link.UserBID, primitive.NilObjectID)
	}

	return nil
}

func (s *UserService) SendMessage(userID, chatroomID primitive.ObjectID, content string) (*model.Message, error) {
	chatroom, err := s.chatroomRepo.GetChatroom(chatroomID)
	if err != nil {
		return nil, err
	}

	if chatroom.UserAID != userID && chatroom.UserBID != userID {
		return nil, errors.New("user is not part of this chatroom")
	}

	if chatroom.IsLocked {
		messages, err := s.chatroomRepo.GetMessages(chatroomID)
		if err != nil {
			return nil, err
		}
		if len(messages) >= 2 {
			return nil, errors.New("chatroom is locked and message limit reached")
		}
	}

	return s.chatroomRepo.AddMessage(chatroomID, userID, content)
}

func (s *UserService) GetMessages(userID, chatroomID primitive.ObjectID) ([]*model.Message, error) {
	chatroom, err := s.chatroomRepo.GetChatroom(chatroomID)
	if err != nil {
		return nil, err
	}

	if chatroom.UserAID != userID && chatroom.UserBID != userID {
		return nil, errors.New("user is not part of this chatroom")
	}

	return s.chatroomRepo.GetMessages(chatroomID)
}

func (s *UserService) UnlockChatroom(chatroomID primitive.ObjectID) error {
	return s.chatroomRepo.UnlockChatroom(chatroomID)
}

func (s *UserService) VerifyNFCAndUnlockChatroom(userID, chatroomID primitive.ObjectID) error {
	// Get the chatroom
	chatroom, err := s.chatroomRepo.GetChatroom(chatroomID)
	if err != nil {
		return err
	}

	// Check if the user is part of this chatroom
	if chatroom.UserAID != userID && chatroom.UserBID != userID {
		return errors.New("user is not part of this chatroom")
	}

	// Check if the chatroom is already unlocked
	if !chatroom.IsLocked {
		return errors.New("chatroom is already unlocked")
	}

	// In a real-world scenario, we might want to implement additional verification here,
	// such as checking if both users have reported an NFC connection within a short time frame.

	// Unlock the chatroom
	err = s.chatroomRepo.UnlockChatroom(chatroomID)
	if err != nil {
		return err
	}

	// You might want to add additional logic here, such as notifying both users that the chatroom is unlocked

	return nil
}
