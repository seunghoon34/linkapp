package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/seunghoon34/linkapp/backend/internal/model"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	collection := db.Collection("users")

	// Create a 2dsphere index on the location field
	_, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "location", Value: "2dsphere"}},
		},
	)
	if err != nil {
		log.Fatalf("Error creating geospatial index: %v", err)
	}

	return &UserRepository{collection: collection}
}

func (r *UserRepository) Create(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)

	return err
}

func (r *UserRepository) UpdateProfile(id string, profile model.Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"profile":   profile,
				"updatedAt": time.Now(),
			},
		},
	)

	return err
}

func (r *UserRepository) UpdatePreferences(id string, preferences model.Preferences) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"preferences": preferences,
				"updatedAt":   time.Now(),
			},
		},
	)

	return err
}

var ErrUserNotFound = errors.New("user not found")

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

const FixedSearchDistance = 200

func (r *UserRepository) SearchMatches(user *model.User, limit int) ([]*model.User, error) {
	ctx := context.Background()

	// Calculate min and max birth dates based on age preferences
	minBirthDate := time.Now().AddDate(-user.Preferences.MaxAge-1, 0, 0)
	maxBirthDate := time.Now().AddDate(-user.Preferences.MinAge, 0, 0)

	pipeline := mongo.Pipeline{
		// Match users based on gender preference, age range, and active search status
		{{Key: "$match", Value: bson.M{
			"is_searching":   true,
			"profile.gender": bson.M{"$in": user.Preferences.Gender},
			"profile.date_of_birth": bson.M{
				"$gte": minBirthDate,
				"$lte": maxBirthDate,
			},
			"_id": bson.M{"$ne": user.ID},
		}}},
		// Calculate distance and filter based on fixed distance
		{{Key: "$geoNear", Value: bson.M{
			"near":          user.Location,
			"distanceField": "distance",
			"maxDistance":   FixedSearchDistance,
			"spherical":     true,
		}}},
		// Check if the current user matches the potential match's preferences
		{{Key: "$match", Value: bson.M{
			"preferences.gender":  user.Profile.Gender,
			"preferences.min_age": bson.M{"$lte": age(user.Profile.DateOfBirth)},
			"preferences.max_age": bson.M{"$gte": age(user.Profile.DateOfBirth)},
		}}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var matches []*model.User
	if err = cursor.All(ctx, &matches); err != nil {
		return nil, err
	}

	return matches, nil
}

func age(birthDate time.Time) int {
	return int(time.Since(birthDate).Hours() / 24 / 365)
}

func (r *UserRepository) UpdateLocation(userID string, latitude, longitude float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"location": model.GeoLocation{
				Type:        "Point",
				Coordinates: []float64{longitude, latitude},
			},
			"updatedAt": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *UserRepository) SetSearchingStatus(userID primitive.ObjectID, isSearching bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"is_searching": isSearching,
			"updated_at":   time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	return err
}

func (r *UserRepository) SetCurrentLink(userID, linkID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"current_link_id": linkID,
			"is_searching":    false,
			"updated_at":      time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	return err
}

func (r *UserRepository) FindPotentialMatch(user *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	minBirthDate := time.Now().AddDate(-user.Preferences.MaxAge-1, 0, 0)
	maxBirthDate := time.Now().AddDate(-user.Preferences.MinAge, 0, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$geoNear", Value: bson.M{
			"near":               user.Location,
			"distanceField":      "distance",
			"maxDistance":        200,
			"spherical":          true,
			"distanceMultiplier": 0.001,
		}}},
		{{Key: "$match", Value: bson.M{
			"_id":            bson.M{"$ne": user.ID},
			"is_searching":   true,
			"profile.gender": bson.M{"$in": user.Preferences.Gender},
			"profile.date_of_birth": bson.M{
				"$gte": minBirthDate,
				"$lte": maxBirthDate,
			},
			"preferences.gender":  user.Profile.Gender,
			"preferences.min_age": bson.M{"$lte": age(user.Profile.DateOfBirth)},
			"preferences.max_age": bson.M{"$gte": age(user.Profile.DateOfBirth)},
		}}},
		{{Key: "$sample", Value: bson.M{"size": 1}}},
	}

	var potentialMatches []model.User
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &potentialMatches); err != nil {
		return nil, err
	}

	if len(potentialMatches) == 0 {
		return nil, nil // No match found
	}

	return &potentialMatches[0], nil
}
