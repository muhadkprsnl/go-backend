package controllers

import (
	"context"
	"net/http"
	"time"

	"local/qa-report/internal/models"
	"local/qa-report/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	UserCollection *mongo.Collection
}

func NewAuthController(client *mongo.Client) *AuthController {
	return &AuthController{
		UserCollection: client.Database("QA").Collection("users"),
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
// 	utils.EnableCORS(w)
// 	if r.Method == http.MethodOptions {
// 		utils.HandleOptions(w, r)
// 		return
// 	}

// 	var req LoginRequest
// 	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
// 		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
// 		return
// 	}

// 	var user models.User
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	err := ac.UserCollection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
// 	if err != nil {
// 		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
// 		return
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
// 		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
// 		return
// 	}

//		token, _ := utils.GenerateJWT(user.Username)
//		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
//	}
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(w)
	if r.Method == http.MethodOptions {
		utils.HandleOptions(w, r)
		return
	}

	var req LoginRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ac.UserCollection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{
			"error":  "Invalid credentials",
			"field":  "username",
			"detail": "User not found",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{
			"error":  "Invalid credentials",
			"field":  "password",
			"detail": "Incorrect password",
		})
		return
	}

	token, _ := utils.GenerateJWT(user.Username)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
