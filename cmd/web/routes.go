package main

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, makeResponseJSON)

	//dynamicMiddleware := alice.New()

	mux := pat.New()

	// Swagger docs
	mux.Get("/swagger/", httpSwagger.WrapHandler)

	// USER
	mux.Post("/api/users", http.HandlerFunc(app.userHandler.Create))
	mux.Get("/api/users", http.HandlerFunc(app.userHandler.GetAll))
	mux.Get("/api/users/:id", http.HandlerFunc(app.userHandler.GetByID))
	mux.Put("/api/users", http.HandlerFunc(app.userHandler.Update))
	mux.Del("/api/users/:id", http.HandlerFunc(app.userHandler.Delete))
	mux.Post("/api/login", http.HandlerFunc(app.userHandler.LoginUser))

	// MICROGREEN
	mux.Post("/api/microgreens", http.HandlerFunc(app.microgreenHandler.Create))
	mux.Get("/api/microgreens", http.HandlerFunc(app.microgreenHandler.GetAll))
	mux.Get("/api/microgreens/:id", http.HandlerFunc(app.microgreenHandler.GetByID))
	mux.Put("/api/microgreens", http.HandlerFunc(app.microgreenHandler.Update))
	mux.Del("/api/microgreens/:id", http.HandlerFunc(app.microgreenHandler.Delete))

	// BATCH
	mux.Post("/api/batches", http.HandlerFunc(app.batchHandler.Create))
	mux.Get("/api/batches", http.HandlerFunc(app.batchHandler.GetAll))
	mux.Get("/api/batches/:id", http.HandlerFunc(app.batchHandler.GetByID))
	mux.Put("/api/batches", http.HandlerFunc(app.batchHandler.Update))
	mux.Del("/api/batches/:id", http.HandlerFunc(app.batchHandler.Delete))

	// OBSERVATION
	mux.Post("/api/observations", http.HandlerFunc(app.observationHandler.Create))
	mux.Get("/api/observations", http.HandlerFunc(app.observationHandler.GetAll))
	mux.Get("/api/observations/:id", http.HandlerFunc(app.observationHandler.GetByID))
	mux.Put("/api/observations", http.HandlerFunc(app.observationHandler.Update))
	mux.Del("/api/observations/:id", http.HandlerFunc(app.observationHandler.Delete))

	// PHOTO
	mux.Post("/api/photos", http.HandlerFunc(app.photoHandler.Create))
	mux.Get("/api/photos", http.HandlerFunc(app.photoHandler.GetAll))
	mux.Get("/api/photos/:id", http.HandlerFunc(app.photoHandler.GetByID))
	mux.Put("/api/photos", http.HandlerFunc(app.photoHandler.Update))
	mux.Del("/api/photos/:id", http.HandlerFunc(app.photoHandler.Delete))
	mux.Get("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// ADVICE
	mux.Post("/api/advice", http.HandlerFunc(app.adviceHandler.Create))
	mux.Get("/api/advice", http.HandlerFunc(app.adviceHandler.GetAll))
	mux.Get("/api/advice/:id", http.HandlerFunc(app.adviceHandler.GetByID))
	mux.Put("/api/advice", http.HandlerFunc(app.adviceHandler.Update))
	mux.Del("/api/advice/:id", http.HandlerFunc(app.adviceHandler.Delete))

	// NOTIFY (FCM)
	// @Tags Notifications
	mux.Post("/notify", http.HandlerFunc(app.fcmHandler.NotifyChange))
	mux.Post("/notify/token/create", http.HandlerFunc(app.fcmHandler.CreateToken))
	mux.Del("/notify/token/:id", http.HandlerFunc(app.fcmHandler.DeleteToken))
	mux.Post("/notify/history", http.HandlerFunc(app.fcmHandler.ShowNotifyHistory))
	mux.Del("/notify/history/:id", http.HandlerFunc(app.fcmHandler.DeleteNotifyHistory))

	// WEBSOCKET
	// @Tags WebSocket
	// @Description Подключение WebSocket для реального времени по пути /ws
	mux.Get("/ws", http.HandlerFunc(app.WebSocketHandler))

	// @Tags WebSocket
	// @Description Подключение WebSocket чата с AI по пути /ws/ai
	mux.Get("/ws/ai", http.HandlerFunc(app.WebSocketAIHandler))

	// REMINDERS
	mux.Post("/api/reminders", http.HandlerFunc(app.reminderHandler.CreateReminder))
	mux.Del("/api/reminders/:id", http.HandlerFunc(app.reminderHandler.DeleteReminder))
	mux.Get("/api/reminders/user/:userId", http.HandlerFunc(app.reminderHandler.GetRemindersByUser))

	//
	mux.Get("/api/analytics/batch-progress", http.HandlerFunc(app.analyticsHandler.GetBatchProgress))
	mux.Get("/api/analytics/humidity-last-7-days", http.HandlerFunc(app.analyticsHandler.GetHumidityLast7Days))
	mux.Get("/api/analytics/height-last-7-days", http.HandlerFunc(app.analyticsHandler.GetHeightLast7Days))
	mux.Get("/api/analytics/ai-prediction", http.HandlerFunc(app.analyticsHandler.GetAIPrediction))
	mux.Get("/api/analytics/today-observations", http.HandlerFunc(app.analyticsHandler.GetTodayObservations))

	return standardMiddleware.Then(mux)
}
