// Для генерации: `swag init`
// Swagger UI:    `localhost:4000/swagger/index.html`

// @title MicroGreens API
// @version 1.0
// @description Backend для проекта "MicroGreens"
// @contact.name Dev Support
// @host localhost:4000
// @BasePath /

package main

import (
	_ "MicroGreens/docs"
	"MicroGreens/internal/config"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
	_ "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	port := os.Getenv("PORT")
	if port != "" {
		port = ":" + port
	} else {
		port = ":4000"
	}

	addr := flag.String("addr", port, "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.Database.URL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	app := initializeApp(db, errorLog, infoLog)

	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"},
	})

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      addSecurityHeaders(c.Handler(app.routes())),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Printf("Starting server on %s", *addr)
	select {}
}
