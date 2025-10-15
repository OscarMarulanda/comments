package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"strconv"
	

	"github.com/OscarMarulanda/comments/internal/database"
	"github.com/OscarMarulanda/comments/internal/models"
	"github.com/OscarMarulanda/comments/internal/utils"
)

/// GET /api/notes
func GetNotes(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	query := r.URL.Query().Get("query")

	var rows *sql.Rows
	var err error

	if query != "" {
		q := "%" + query + "%"
		rows, err = database.DB.Query(
			`SELECT id, title, content, user_id, created_at, updated_at, deleted_at
			 FROM notes
			 WHERE user_id = $1 
			 AND deleted_at IS NULL
			 AND (title ILIKE $2 OR content ILIKE $2)
			 ORDER BY id DESC`,
			userID, q,
		)
	} else {
		rows, err = database.DB.Query(
			`SELECT id, title, content, user_id, created_at, updated_at, deleted_at
			 FROM notes
			 WHERE user_id = $1 
			 AND deleted_at IS NULL
			 ORDER BY id DESC`,
			userID,
		)
	}

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch notes")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(
			&n.ID, &n.Title, &n.Content, &n.UserID,
			&n.CreatedAt, &n.UpdatedAt, &n.DeletedAt,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notes = append(notes, n)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// POST /api/notes
func CreateNote(w http.ResponseWriter, r *http.Request) {
	var n models.Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		log.Println("Decode error:", err)
		return
	}

	userIDVal := r.Context().Value("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		log.Println("No userID in context or wrong type")
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := database.DB.QueryRow(
		"INSERT INTO notes (title, content, user_id) VALUES ($1, $2, $3) RETURNING id, created_at",
		n.Title, n.Content, userID,
	).Scan(&n.ID, &n.CreatedAt)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create note")
		log.Println("DB insert error:", err)
		return
	}

	n.UserID = userID
	utils.WriteSuccess(w, n)
}

// PUT /api/notes/{id}
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var n models.Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		log.Println("Decode error:", err)
		return
	}

	userID := r.Context().Value("userID").(int)

	// Run the UPDATE query
	_, err := database.DB.Exec(
		`UPDATE notes 
         SET title = $1, content = $2, updated_at = NOW() 
         WHERE id = $3 AND user_id = $4`,
		n.Title, n.Content, id, userID,
	)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update note")
		log.Println("DB update error:", err)
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Note updated successfully"})
}

// DELETE /api/notes/{id}
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	userID := r.Context().Value("userID").(int)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid note ID")
		log.Println("Invalid ID:", idStr)
		return
	}

	result, err := database.DB.Exec("UPDATE notes SET deleted_at = NOW() WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL",
		id, userID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete note")
		log.Println("DB delete error:", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Note not found")
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Note deleted successfully"})
}
