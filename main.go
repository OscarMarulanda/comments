package main

import (
    "github.com/joho/godotenv"
    "log"
    "net/http"
    "github.com/OscarMarulanda/comments/internal/routes"
	"github.com/OscarMarulanda/comments/internal/database"
    "os"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  No .env file found (using system environment variables)")
    }
    
	if err := database.Connect(); err != nil {
        log.Fatal("❌ Database connection failed:", err)
    }
	
    r := routes.SetupRouter()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Println("Server running on :" + port)
    http.ListenAndServe(":"+port, r)
}