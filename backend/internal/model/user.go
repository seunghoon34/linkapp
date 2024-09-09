package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username      string             `bson:"username" json:"username"`
	Email         string             `bson:"email" json:"email"`
	Password      string             `bson:"password" json:"-"`
	Profile       Profile            `bson:"profile" json:"profile"`
	Preferences   Preferences        `bson:"preferences" json:"preferences"`
	Location      GeoLocation        `bson:"location" json:"location"`
	IsSearching   bool               `bson:"is_searching" json:"is_searching"`
	CurrentLinkID primitive.ObjectID `bson:"current_link_id,omitempty" json:"current_link_id,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type GeoLocation struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

type Profile struct {
	FirstName     string    `bson:"first_name" json:"first_name"`
	LastName      string    `bson:"last_name" json:"last_name"`
	DateOfBirth   time.Time `bson:"date_of_birth" json:"date_of_birth"`
	Gender        string    `bson:"gender" json:"gender"`
	Bio           string    `bson:"bio" json:"bio"`
	ProfilePicURL string    `bson:"profile_pic_url" json:"profile_pic_url"`
}

type Preferences struct {
	MinAge int      `bson:"min_age" json:"min_age"`
	MaxAge int      `bson:"max_age" json:"max_age"`
	Gender []string `bson:"gender" json:"gender"`
}

type Link struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserAID   primitive.ObjectID `bson:"user_a_id" json:"user_a_id"`
	UserBID   primitive.ObjectID `bson:"user_b_id" json:"user_b_id"`
	Status    LinkStatus         `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
}

type LinkStatus string

const (
	LinkStatusPending  LinkStatus = "pending"
	LinkStatusAccepted LinkStatus = "accepted"
	LinkStatusRejected LinkStatus = "rejected"
	LinkStatusExpired  LinkStatus = "expired"
)
