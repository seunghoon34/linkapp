package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/seunghoon34/linkapp/backend/internal/model"
	"github.com/seunghoon34/linkapp/backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}
	err := h.userService.CreateUser(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var profile model.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdateProfile(id, profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var preferences model.Preferences
	if err := json.NewDecoder(r.Body).Decode(&preferences); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdatePreferences(id, preferences); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdateUser(id, input.Username, input.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Here you would typically generate a JWT token and return it
	// For now, we'll just return the user object
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
func (h *UserHandler) SearchMatches(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20 // Default limit
	}

	matches, err := h.userService.SearchMatches(userID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

// You can add more handler methods as needed

func (h *UserHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var locationData struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := json.NewDecoder(r.Body).Decode(&locationData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.userService.UpdateLocation(userID, locationData.Latitude, locationData.Longitude)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Location updated successfully"})
}

func (h *UserHandler) StartSearching(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.userService.StartSearching(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Started searching"})
}

func (h *UserHandler) StopSearching(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.userService.StopSearching(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Stopped searching"})
}

func (h *UserHandler) FindMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	link, err := h.userService.FindMatch(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(link)
}

func (h *UserHandler) RespondToLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	linkID, err := primitive.ObjectIDFromHex(vars["linkId"])
	if err != nil {
		http.Error(w, "Invalid link ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Accept bool `json:"accept"`
	}
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.userService.RespondToLink(userID, linkID, response.Accept)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Response recorded"})
}

func (h *UserHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	chatroomID, err := primitive.ObjectIDFromHex(vars["chatroomId"])
	if err != nil {
		http.Error(w, "Invalid chatroom ID", http.StatusBadRequest)
		return
	}

	var messageData struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&messageData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := h.userService.SendMessage(userID, chatroomID, messageData.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (h *UserHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	chatroomID, err := primitive.ObjectIDFromHex(vars["chatroomId"])
	if err != nil {
		http.Error(w, "Invalid chatroom ID", http.StatusBadRequest)
		return
	}

	messages, err := h.userService.GetMessages(userID, chatroomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

func (h *UserHandler) UnlockChatroom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatroomID, err := primitive.ObjectIDFromHex(vars["chatroomId"])
	if err != nil {
		http.Error(w, "Invalid chatroom ID", http.StatusBadRequest)
		return
	}

	err = h.userService.UnlockChatroom(chatroomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Chatroom unlocked"})
}

func (h *UserHandler) VerifyNFCAndUnlockChatroom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	chatroomID, err := primitive.ObjectIDFromHex(vars["chatroomId"])
	if err != nil {
		http.Error(w, "Invalid chatroom ID", http.StatusBadRequest)
		return
	}

	err = h.userService.VerifyNFCAndUnlockChatroom(userID, chatroomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Chatroom unlocked via NFC"})
}
