// handlers/user.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"task-project/models"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateMyUserDummyTesting(userCol *mongo.Collection, profileCol *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		// Decode JSON dari Postman
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
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
		_, err := userCol.InsertOne(r.Context(), user)
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
