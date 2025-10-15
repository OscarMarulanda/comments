package handlers

import (
    
    "net/http"

    "github.com/OscarMarulanda/comments/internal/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
    utils.WriteSuccess(w, map[string]string{"status": "ok"})
}
