package repositories

import (
	"context"
	"local/qa-report/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	collection *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) *AuthRepository {
	return &AuthRepository{
		collection: db.Collection("users"),
	}
}

func (r *AuthRepository) CreateUser(user *models.User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	_, err = r.collection.InsertOne(context.Background(), user)
	return err
}

func (r *AuthRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
