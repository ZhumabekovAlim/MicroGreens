package main

import (
	"MicroGreens/internal/handlers"
	"context"
	"database/sql"
	"log"
	"time"
)

type Reminder struct {
	ID      int
	UserID  int
	Message string
	Time    string
	Active  bool
}

func startReminderScheduler(db *sql.DB, fcmHandler *handlers.FCMHandler) {
	go func() {
		for {
			now := time.Now().Format("15:04")

			reminders, err := getDueReminders(db, now)
			if err != nil {
				log.Println("Ошибка при получении напоминаний:", err)
			}

			for _, r := range reminders {
				tokens, err := fcmHandler.GetTokensByClientID(r.UserID)
				if err != nil {
					log.Printf("Ошибка токенов для user_id %d: %v", r.UserID, err)
					continue
				}

				for _, token := range tokens {
					err = fcmHandler.SendMessage(context.Background(), token, r.UserID, 0, r.UserID, "Ежедневное напоминание", r.Message)
					if err != nil {
						log.Printf("Ошибка при отправке напоминания: %v", err)
					}
				}
			}

			time.Sleep(10 * time.Second)
		}
	}()
}

func getDueReminders(db *sql.DB, currentTime string) ([]Reminder, error) {
	rows, err := db.Query(`SELECT id, user_id, message, time, active FROM reminders WHERE active = TRUE AND time = ?`, currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []Reminder
	for rows.Next() {
		var r Reminder
		if err := rows.Scan(&r.ID, &r.UserID, &r.Message, &r.Time, &r.Active); err != nil {
			return nil, err
		}
		reminders = append(reminders, r)
	}
	return reminders, nil
}
