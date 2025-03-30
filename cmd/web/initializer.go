package main

import (
	"MicroGreens/internal/handlers"
	"MicroGreens/internal/repositories"
	"MicroGreens/internal/services"
	"context"
	"database/sql"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

type application struct {
	db                 *sql.DB
	errorLog           *log.Logger
	infoLog            *log.Logger
	wsManager          *WebSocketManager
	microgreenHandler  *handlers.MicrogreenHandler
	fcmHandler         *handlers.FCMHandler
	batchHandler       *handlers.BatchHandler
	observationHandler *handlers.ObservationHandler
	photoHandler       *handlers.PhotoHandler
	adviceHandler      *handlers.AdviceHandler
	userHandler        *handlers.UserHandler
	reminderHandler    *handlers.ReminderHandler
	analyticsHandler   *handlers.AnalyticsHandler
}

func initializeApp(db *sql.DB, errorLog, infoLog *log.Logger) *application {

	ctx := context.Background()
	sa := option.WithCredentialsFile("C:\\Users\\alimz\\GolandProjects\\MicroGreens\\cmd\\web\\serviceAccountKey.json")
	if _, err := os.Stat("C:\\Users\\alimz\\GolandProjects\\MicroGreens\\cmd\\web\\serviceAccountKey.json"); err != nil {
		log.Fatalf("serviceAccountKey.json not found: %v", err)
	}

	firebaseApp, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: "greending-470ce"}, sa)
	if err != nil {
		errorLog.Fatalf("Ошибка в нахождении приложения: %v\n", err)
	}

	fcmClient, err := firebaseApp.Messaging(ctx)
	if err != nil {
		errorLog.Fatalf("Ошибка при неверном ID устройства: %v\n", err)
	}

	fcmHandler := handlers.NewFCMHandler(fcmClient, db)

	userRepo := &repositories.UserRepository{Db: db}
	userService := &services.UserService{Repo: userRepo}
	userHandler := &handlers.UserHandler{Service: userService}

	microgreenRepo := &repositories.MicrogreenRepository{Db: db}
	microgreenService := &services.MicrogreenService{Repo: microgreenRepo}
	microgreenHandler := &handlers.MicrogreenHandler{Service: microgreenService}

	batchRepo := &repositories.BatchRepository{Db: db}
	batchService := &services.BatchService{Repo: batchRepo}
	batchHandler := &handlers.BatchHandler{Service: batchService}

	observationRepo := &repositories.ObservationRepository{Db: db}
	observationService := &services.ObservationService{Repo: observationRepo}
	observationHandler := &handlers.ObservationHandler{Service: observationService}

	photoRepo := &repositories.PhotoRepository{Db: db}
	photoService := &services.PhotoService{Repo: photoRepo}
	photoHandler := &handlers.PhotoHandler{Service: photoService}

	adviceRepo := &repositories.AdviceRepository{Db: db}
	adviceService := &services.AdviceService{Repo: adviceRepo}
	adviceHandler := &handlers.AdviceHandler{Service: adviceService}

	reminderHandler := handlers.NewReminderHandler(db)

	analyticsHandler := handlers.NewAnalyticsHandler(db)

	return &application{
		db:                 db,
		errorLog:           errorLog,
		infoLog:            infoLog,
		wsManager:          NewWebSocketManager(),
		fcmHandler:         fcmHandler,
		microgreenHandler:  microgreenHandler,
		batchHandler:       batchHandler,
		observationHandler: observationHandler,
		photoHandler:       photoHandler,
		adviceHandler:      adviceHandler,
		userHandler:        userHandler,
		reminderHandler:    reminderHandler,
		analyticsHandler:   analyticsHandler,
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("%v", err)
		panic("failed to connect to database")
		return nil, err
	}
	db.SetMaxIdleConns(35)
	if err = db.Ping(); err != nil {
		log.Printf("%v", err)
		panic("failed to ping the database")
		return nil, err
	}
	fmt.Println("successfully connected")

	return db, nil
}

func addSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}
