package repository

import (
	"context"
	"time"

	"github.com/seunghoon34/linkapp/backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatroomRepository struct {
	chatroomCollection *mongo.Collection
	messageCollection  *mongo.Collection
}

func NewChatroomRepository(db *mongo.Database) *ChatroomRepository {
	return &ChatroomRepository{
		chatroomCollection: db.Collection("chatrooms"),
		messageCollection:  db.Collection("messages"),
	}
}

func (r *ChatroomRepository) CreateChatroom(linkID, userAID, userBID primitive.ObjectID) (*model.Chatroom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom := &model.Chatroom{
		LinkID:    linkID,
		UserAID:   userAID,
		UserBID:   userBID,
		IsLocked:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := r.chatroomCollection.InsertOne(ctx, chatroom)
	if err != nil {
		return nil, err
	}

	chatroom.ID = result.InsertedID.(primitive.ObjectID)
	return chatroom, nil
}

func (r *ChatroomRepository) GetChatroom(chatroomID primitive.ObjectID) (*model.Chatroom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var chatroom model.Chatroom
	err := r.chatroomCollection.FindOne(ctx, bson.M{"_id": chatroomID}).Decode(&chatroom)
	if err != nil {
		return nil, err
	}

	return &chatroom, nil
}

func (r *ChatroomRepository) UnlockChatroom(chatroomID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"is_locked":  false,
			"updated_at": time.Now(),
		},
	}

	_, err := r.chatroomCollection.UpdateOne(ctx, bson.M{"_id": chatroomID}, update)
	return err
}

func (r *ChatroomRepository) AddMessage(chatroomID, senderID primitive.ObjectID, content string) (*model.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message := &model.Message{
		ChatroomID: chatroomID,
		SenderID:   senderID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	result, err := r.messageCollection.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}

	message.ID = result.InsertedID.(primitive.ObjectID)
	return message, nil
}

func (r *ChatroomRepository) GetMessages(chatroomID primitive.ObjectID) ([]*model.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.messageCollection.Find(ctx, bson.M{"chatroom_id": chatroomID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*model.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
