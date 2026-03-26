package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func CheckUserAvailability(email string) bool {
	collection := GetCollection("users")
	var result bson.M
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&result)
	return err == mongo.ErrNoDocuments
}

func UserCreate(email string, password string) *User {
	// Password Hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}

	user := User{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := GetCollection("users")
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return nil
	}

	return &user
}

func UserMatchPassword(email string, password string) *User {
	collection := GetCollection("users")
	var user User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return &User{}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return &User{}
	} else {
		return &user
	}
}

func UserFromId(id primitive.ObjectID) *User {
	collection := GetCollection("users")
	var user User
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return &User{}
	}
	return &user
}
