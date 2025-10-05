package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"task-project/models"
)

func CreateTask(taskCol *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		// decode json from postmand
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// get user id from token authmiddleware
		userIDHex := r.Header.Get("X-User-ID")
		userID, err := primitive.ObjectIDFromHex(userIDHex)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		if task.Status == "" {
			task.Status = "pending"
		}

		now := time.Now()
		task.ID = primitive.NewObjectID()
		task.UserID = userID
		task.CreatedAt = now
		task.UpdatedAt = now

		task.StatusHistory = []models.StatusChange{
			{
				From:      "",
				To:        task.Status,
				ChangedAt: now,
				ChangedBy: userID,
			},
		}

		// Insert to task collection
		_, err = taskCol.InsertOne(r.Context(), task)
		if err != nil {
			log.Println("❌ Failed to insert task:", err)
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
			return
		}

		log.Println("✅ Task created:", task.Title)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Task created successfully"))
	}
}

func GetTaskByID(taskCol *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := strings.TrimSpace(r.URL.Query().Get("taskID"))

		var filter bson.M
		if oid, err := primitive.ObjectIDFromHex(idParam); err == nil {
			filter = bson.M{"_id": oid}
		} else {
			filter = bson.M{"_id": idParam}
		}

		var task models.Task
		err := taskCol.FindOne(r.Context(), filter).Decode(&task)
		if err != nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

func GetAllTasks(taskCol *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query params
		status := r.URL.Query().Get("status")
		priority := r.URL.Query().Get("priority")
		title := r.URL.Query().Get("title")
		sortField := r.URL.Query().Get("sort")
		order := r.URL.Query().Get("order")
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		// Default pagination
		page, _ := strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(limitStr)
		if limit < 1 {
			limit = 10
		}
		skip := (page - 1) * limit

		// Build filter
		filter := bson.M{}
		if status != "" && strings.ToLower(status) != "all" {
			filter["status"] = status
		}
		if priority != "" && strings.ToLower(priority) != "all" {
			filter["priority"] = priority
		}
		if title != "" {
			filter["title"] = bson.M{"$regex": title, "$options": "i"}
		}

		// Build sort
		sortOrder := -1 // default descending
		if strings.ToLower(order) == "asc" {
			sortOrder = 1
		}
		sortBy := bson.D{{Key: "due_date", Value: sortOrder}} // default sort by due_date
		if sortField != "" {
			sortBy = bson.D{{Key: sortField, Value: sortOrder}}
		}

		// Count total
		total, err := taskCol.CountDocuments(r.Context(), filter)
		if err != nil {
			http.Error(w, "Failed to count tasks", http.StatusInternalServerError)
			return
		}

		// Fetch data
		cursor, err := taskCol.Find(r.Context(), filter, &options.FindOptions{
			Sort:  sortBy,
			Skip:  int64Ptr(int64(skip)),
			Limit: int64Ptr(int64(limit)),
		})
		if err != nil {
			http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		var tasks []models.Task
		if err := cursor.All(r.Context(), &tasks); err != nil {
			http.Error(w, "Failed to decode tasks", http.StatusInternalServerError)
			return
		}

		// Response
		resp := map[string]interface{}{
			"data": tasks,
			"meta": map[string]interface{}{
				"total": total,
				"page":  page,
				"limit": limit,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func int64Ptr(i int64) *int64 {
	return &i
}
