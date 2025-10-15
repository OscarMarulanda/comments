package utils

import (
    "encoding/json"
    "net/http"
)

// JSONResponse represents a consistent API response format.
type JSONResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// WriteJSON sends a JSON response with a given status code and payload.
func WriteJSON(w http.ResponseWriter, status int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(payload)
}

// WriteError sends a standardized error response.
func WriteError(w http.ResponseWriter, status int, message string) {
    resp := JSONResponse{
        Success: false,
        Error:   message,
    }
    WriteJSON(w, status, resp)
}

// WriteSuccess sends a standardized success response.
func WriteSuccess(w http.ResponseWriter, data interface{}) {
    resp := JSONResponse{
        Success: true,
        Data:    data,
    }
    WriteJSON(w, http.StatusOK, resp)
}