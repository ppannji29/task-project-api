package routes

import (
	"net/http"
	"task-project/db"
	"task-project/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(store *db.MongoCollections) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test mock API is running!"))
	})

	// User dummy creation
	r.Post("/api/user/mydummy/create", handlers.CreateMyUserDummyTesting(store.UserCol, store.ProfileCol))

	// Auth
	r.Post("/api/auth/request-otp", handlers.RequestOtp(store.UserCol, store.OtpCol))
	r.Post("/api/auth/verify-otp", handlers.VerifyOtp(store.UserCol, store.OtpCol))
	r.Post("/api/auth/refresh-token", handlers.RefreshToken(store.UserCol))
	r.Post("/api/auth/logout", handlers.Logout())

	// Tasks (protected)
	// r.Group(func(r chi.Router) {
	//     r.Use(handlers.AuthMiddleware)
	//     r.Get("/api/tasks", handlers.GetTasks())
	//     r.Post("/api/tasks", handlers.CreateTask())
	//     r.Put("/api/tasks/{id}", handlers.UpdateTask())
	//     r.Delete("/api/tasks/{id}", handlers.DeleteTask())
	// })

	return r
}

// package routes

// import (
// 	"net/http"
// 	"task-project/db"
// 	"task-project/handlers"
// )

// func SetupRoutes(store *db.MongoCollections) *http.ServeMux {
// 	mux := http.NewServeMux()
// 	// api check health
// 	mux.HandleFunc("/", healthCheck)
// 	// create dummy user for testing
// 	mux.HandleFunc("POST /api/user/mydummy/create", handlers.CreateMyUserDummyTesting(store.UserCol, store.ProfileCol))
// 	// Auth
// 	mux.HandleFunc("POST /api/auth/request-otp", handlers.RequestOtp(store.UserCol, store.OtpCol))
// 	mux.HandleFunc("POST /api/auth/verify-otp", handlers.VerifyOtp(store.UserCol, store.OtpCol))
// 	mux.HandleFunc("POST /api/auth/refresh-token", handlers.RefreshToken(store.UserCol))
// 	mux.HandleFunc("POST /api/auth/logout", handlers.Logout())

// 	// // Task CRUD (protected)
// 	// mux.Handle("GET /api/tasks", handlers.AuthMiddleware(http.HandlerFunc(handlers.GetTasks())))
// 	// mux.Handle("POST /api/tasks", handlers.AuthMiddleware(http.HandlerFunc(handlers.CreateTask())))
// 	// mux.Handle("PUT /api/tasks/", handlers.AuthMiddleware(http.HandlerFunc(handlers.UpdateTask()))) // use query param or path parsing
// 	// mux.Handle("DELETE /api/tasks/", handlers.AuthMiddleware(http.HandlerFunc(handlers.DeleteTask())))

// 	return mux
// }

// func healthCheck(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Test mock API is running!"))
// }
