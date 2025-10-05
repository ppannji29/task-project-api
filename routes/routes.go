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
	r.Delete("/api/auth/logout", handlers.Logout())

	// Tasks (protected)
	r.Route("/api", func(r chi.Router) {
		r.Use(handlers.AuthMiddleware())

		r.Get("/user/me", handlers.GetCurrentUser(store.UserCol))
		r.Get("/tasks", handlers.GetAllTasks(store.TaskCol))
		r.Get("/task/{id}", handlers.GetTaskByID(store.TaskCol, store.UserCol))
		r.Post("/task", handlers.CreateTask(store.TaskCol))
		r.Patch("/task/{id}", handlers.UpdateTask(store.TaskCol))
		r.Delete("/task/{id}", handlers.DeleteTask(store.TaskCol))
	})

	return r
}
