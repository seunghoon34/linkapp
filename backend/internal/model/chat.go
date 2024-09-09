package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chatroom struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	LinkID    primitive.ObjectID `bson:"link_id" json:"link_id"`
	UserAID   primitive.ObjectID `bson:"user_a_id" json:"user_a_id"`
	UserBID   primitive.ObjectID `bson:"user_b_id" json:"user_b_id"`
	IsLocked  bool               `bson:"is_locked" json:"is_locked"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Message struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ChatroomID primitive.ObjectID `bson:"chatroom_id" json:"chatroom_id"`
	SenderID   primitive.ObjectID `bson:"sender_id" json:"sender_id"`
	Content    string             `bson:"content" json:"content"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
