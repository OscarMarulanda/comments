package database

import (
    "database/sql"
    "fmt"
    "os"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        return fmt.Errorf("DATABASE_URL is not set")
    }

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return err
    }

    if err := db.Ping(); err != nil {
        return err
    }

    DB = db
    fmt.Println("✅ Connected to PostgreSQL")
    return nil
}