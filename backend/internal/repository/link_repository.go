package repository

import (
	"context"
	"time"

	"github.com/seunghoon34/linkapp/backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LinkRepository struct {
	collection *mongo.Collection
}

func NewLinkRepository(db *mongo.Database) *LinkRepository {
	return &LinkRepository{
		collection: db.Collection("links"),
	}
}

func (r *LinkRepository) CreateLink(userAID, userBID primitive.ObjectID) (*model.Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	link := &model.Link{
		UserAID:   userAID,
		UserBID:   userBID,
		Status:    model.LinkStatusPending,
		CreatedAt: now,
		ExpiresAt: now.Add(30 * time.Second),
	}

	result, err := r.collection.InsertOne(ctx, link)
	if err != nil {
		return nil, err
	}

	link.ID = result.InsertedID.(primitive.ObjectID)
	return link, nil
}

func (r *LinkRepository) UpdateLinkStatus(linkID primitive.ObjectID, status model.LinkStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": linkID}, update)
	return err
}

func (r *LinkRepository) GetLink(linkID primitive.ObjectID) (*model.Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var link model.Link
	err := r.collection.FindOne(ctx, bson.M{"_id": linkID}).Decode(&link)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *LinkRepository) ExpireLinks() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{
		"status":     model.LinkStatusPending,
		"expires_at": bson.M{"$lt": now},
	}
	update := bson.M{
		"$set": bson.M{"status": model.LinkStatusExpired},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *LinkRepository) GetExpiredLinks() ([]*model.Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{
		"status":     model.LinkStatusExpired,
		"expires_at": bson.M{"$lt": now},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var links []*model.Link
	if err = cursor.All(ctx, &links); err != nil {
		return nil, err
	}

	return links, nil
}
