package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Blog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func BlogsAll() *[]Blog {
	collection := GetCollection("blogs")
	var blogs []Blog
	opts := options.Find().SetSort(bson.M{"updated_at": -1})
	cursor, err := collection.Find(context.Background(), bson.M{"deleted_at": nil}, opts)
	if err != nil {
		return &blogs
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &blogs); err != nil {
		return &blogs
	}
	return &blogs
}

func BlogsFind(id primitive.ObjectID) *Blog {
	collection := GetCollection("blogs")
	var blog Blog
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&blog)
	if err != nil {
		return nil
	}
	return &blog
}
