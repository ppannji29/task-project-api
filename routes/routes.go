package routes

import (
	"net/http"
	"task-project/db"
	"task-project/handlers"
)

func SetupRoutes(store *db.MongoCollections) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", healthCheck)

	// Auth
	mux.HandleFunc("POST /api/auth/request-otp", handlers.RequestOtp(store.UserCol, store.OtpCol))
	mux.HandleFunc("POST /api/auth/verify-otp", handlers.VerifyOtp(store.UserCol, store.OtpCol))
	mux.HandleFunc("POST /api/auth/refresh-token", handlers.RefreshToken(store.UserCol))
	mux.HandleFunc("POST /api/auth/logout", handlers.Logout())

	// // Task CRUD (protected)
	// mux.Handle("GET /api/tasks", handlers.AuthMiddleware(http.HandlerFunc(handlers.GetTasks())))
	// mux.Handle("POST /api/tasks", handlers.AuthMiddleware(http.HandlerFunc(handlers.CreateTask())))
	// mux.Handle("PUT /api/tasks/", handlers.AuthMiddleware(http.HandlerFunc(handlers.UpdateTask()))) // use query param or path parsing
	// mux.Handle("DELETE /api/tasks/", handlers.AuthMiddleware(http.HandlerFunc(handlers.DeleteTask())))

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("MileApp mock API is running!"))
}
