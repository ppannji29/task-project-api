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

	"task-project/dto"
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

func GetTaskByID(taskCol *mongo.Collection, userCol *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := strings.Split(r.URL.Path, "/")
		idParam := urlPath[len(urlPath)-1]
		log.Println("✅ ID Param:", idParam)

		if idParam == "" {
			http.Error(w, "Missing task ID in URL", http.StatusBadRequest)
			return
		}

		oid, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			http.Error(w, "Invalid task ID format", http.StatusBadRequest)
			return
		}

		var task models.Task
		err = taskCol.FindOne(r.Context(), bson.M{"_id": oid}).Decode(&task)
		if err != nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		// Enrich status history with email
		for i, history := range task.StatusHistory {
			var user models.User
			err := userCol.FindOne(r.Context(), bson.M{"_id": history.ChangedBy}).Decode(&user)
			if err == nil {
				task.StatusHistory[i].ChangedByEmail = user.Email
			}
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

func UpdateTask(taskCol *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := strings.Split(r.URL.Path, "/")
		idParam := urlPath[len(urlPath)-1]
		taskID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			log.Println("❌ Invalid task ID:", err)
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		var req dto.UpdateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		userIDHex := r.Header.Get("X-User-ID")
		userID, err := primitive.ObjectIDFromHex(userIDHex)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		var existing models.Task
		err = taskCol.FindOne(r.Context(), bson.M{"_id": taskID, "user_id": userID}).Decode(&existing)
		if err != nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		update := bson.M{}
		history := []models.StatusChange{}
		now := time.Now()

		if req.Title != nil {
			update["title"] = *req.Title
		}
		if req.Description != nil {
			update["description"] = *req.Description
		}
		if req.Status != nil && *req.Status != existing.Status {
			update["status"] = *req.Status
			history = append(history, models.StatusChange{
				Type:      "status",
				Action:    "Status Updated",
				From:      existing.Status,
				To:        *req.Status,
				ChangedAt: now,
				ChangedBy: userID,
			})
		}
		if req.Priority != nil && *req.Priority != existing.Priority {
			update["priority"] = *req.Priority
			history = append(history, models.StatusChange{
				Type:      "priority",
				Action:    "Priority Changed",
				From:      existing.Priority,
				To:        *req.Priority,
				ChangedAt: now,
				ChangedBy: userID,
			})
		}
		if req.DueDate != nil && !req.DueDate.Equal(existing.DueDate) {
			update["due_date"] = *req.DueDate
			history = append(history, models.StatusChange{
				Type:      "due_date",
				Action:    "Due Date Rescheduled",
				From:      existing.DueDate.Format("2006-01-02"),
				To:        req.DueDate.Format("2006-01-02"),
				ChangedAt: now,
				ChangedBy: userID,
			})
		}

		update["updated_at"] = now

		ops := bson.M{"$set": update}
		if len(history) > 0 {
			ops["$push"] = bson.M{"status_history": bson.M{"$each": history}}
		}

		_, err = taskCol.UpdateOne(r.Context(), bson.M{"_id": taskID}, ops)
		if err != nil {
			log.Println("❌ Failed to update task:", err)
			http.Error(w, "Failed to update task", http.StatusInternalServerError)
			return
		}

		log.Println("✅ Task updated:", taskID.Hex())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Task updated successfully"))
	}
}

func DeleteTask(taskCol *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := strings.Split(r.URL.Path, "/")
		idParam := urlPath[len(urlPath)-1]
		oid, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		_, err = taskCol.DeleteOne(r.Context(), bson.M{"_id": oid})
		if err != nil {
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
