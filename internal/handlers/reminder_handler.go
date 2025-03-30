package handlers

import (
	"MicroGreens/internal/models"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type ReminderHandler struct {
	DB *sql.DB
}

func NewReminderHandler(db *sql.DB) *ReminderHandler {
	return &ReminderHandler{DB: db}
}

func (h *ReminderHandler) CreateReminder(w http.ResponseWriter, r *http.Request) {
	var reminder models.Reminder
	if err := json.NewDecoder(r.Body).Decode(&reminder); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO reminders (user_id, message, time, active) VALUES (?, ?, ?, ?)`
	_, err := h.DB.Exec(query, reminder.UserID, reminder.Message, reminder.Time, true)
	if err != nil {
		log.Println("DB error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ReminderHandler) DeleteReminder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	query := `DELETE FROM reminders WHERE id = ?`
	_, err := h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete reminder", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ReminderHandler) GetRemindersByUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userId"]
	query := `SELECT id, user_id, message, time, active FROM reminders WHERE user_id = ?`

	rows, err := h.DB.Query(query, userId)
	if err != nil {
		http.Error(w, "Failed to fetch reminders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reminders []models.Reminder
	for rows.Next() {
		var r models.Reminder
		if err := rows.Scan(&r.ID, &r.UserID, &r.Message, &r.Time, &r.Active); err != nil {
			http.Error(w, "Error reading reminder", http.StatusInternalServerError)
			return
		}
		reminders = append(reminders, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}
