// handlers/user.go
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"task-project/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateMyUserDummyTesting(userCol *mongo.Collection, profileCol *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if user.Email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		count, err := userCol.CountDocuments(r.Context(), bson.M{"email": user.Email})
		if err != nil {
			http.Error(w, "Failed to check email", http.StatusInternalServerError)
			return
		}
		if count > 0 {
			http.Error(w, "Your account is already registered", http.StatusConflict)
			return
		}

		// Set default values
		now := time.Now()
		user.CreatedAt = now
		user.UpdatedAt = now
		user.IsEnabled = true

		if user.Profile != nil {
			user.Profile.CreatedAt = now
			user.Profile.UpdatedAt = now
		}

		// Insert user
		_, err = userCol.InsertOne(r.Context(), user)
		if err != nil {
			log.Println("❌ Failed to insert user:", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		log.Println("✅ User created:", user.Email)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created successfully"))
	}
}

func GetCurrentUser(userCol *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ambil user ID dari header (AuthMiddleware lo udah set ini)
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Convert ke ObjectID
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Query user dari Mongo
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.User
		err = userCol.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Profile retrieved successfully",
			"data":    user,
		})
	}
}
