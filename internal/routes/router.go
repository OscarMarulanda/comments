package routes

import (
	"github.com/OscarMarulanda/comments/internal/handlers"
    "github.com/OscarMarulanda/comments/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// CORS middleware
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(commonHeaders)

    // Public routes
	r.HandleFunc("/api/health", handlers.HealthCheck).Methods("GET")
	r.HandleFunc("/api/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")

	// Protected routes — require JWT
	notes := r.PathPrefix("/api/notes").Subrouter()
	notes.Use(middleware.AuthMiddleware)

	notes.HandleFunc("", handlers.GetNotes).Methods("GET")
	notes.HandleFunc("", handlers.CreateNote).Methods("POST")
	notes.HandleFunc("/{id}", handlers.UpdateNote).Methods("PUT")
	notes.HandleFunc("/{id}", handlers.DeleteNote).Methods("DELETE")

	return r
}

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		next.ServeHTTP(w, r)
	})
}
