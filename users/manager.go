package users

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// User represents a user in the system
type User struct {
	UserID    string    `bson:"user_id"`
	Email     string    `bson:"email"`
	Name      string    `bson:"name"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// UserManager handles user operations
type UserManager struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewUserManager creates a new user manager
func NewUserManager(db *mongo.Database) *UserManager {
	return &UserManager{
		db:         db,
		collection: db.Collection("users"),
	}
}

// AddUser adds a new user
func (um *UserManager) AddUser(ctx context.Context, userID, email, name, password string) error {
	// Check if user already exists
	count, err := um.collection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("user already exists")
	}

	// Create new user
	user := User{
		UserID:    userID,
		Email:     email,
		Name:      name,
		Password:  password, // In a real application, this should be hashed
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = um.collection.InsertOne(ctx, user)
	return err
}

// UpdateUser updates an existing user
func (um *UserManager) UpdateUser(ctx context.Context, userID, email, name, password string) error {
	update := bson.M{
		"$set": bson.M{
			"email":      email,
			"name":       name,
			"password":   password, // In a real application, this should be hashed
			"updated_at": time.Now(),
		},
	}

	result, err := um.collection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

// DeleteUser deletes a user
func (um *UserManager) DeleteUser(ctx context.Context, userID string) error {
	result, err := um.collection.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

// AuthenticateUser verifies user credentials
func (um *UserManager) AuthenticateUser(ctx context.Context, userID, password string) (bool, error) {
	var user User
	err := um.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, errors.New("user not found")
		}
		return false, err
	}

	// In a real application, you would compare hashed passwords
	if user.Password != password {
		return false, errors.New("invalid password")
	}

	return true, nil
}