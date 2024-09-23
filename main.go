package main

import (
	"birthday-tracker/internal/database"
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pressly/goose/v3"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

//go:embed static
var static embed.FS

//go:generate npm run --prefix web/tailwindcss build

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL not found!")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(conn, "sql/schema"); err != nil {
		panic(err)
	}

	http.HandleFunc("GET /", apiCfg.handlerGetHomePage)
	http.HandleFunc("POST /", apiCfg.handlerCreatePerson)
	http.HandleFunc("DELETE /{id}", apiCfg.handlerRemovePerson)

	http.Handle("GET /static/", http.FileServer(http.FS(static)))

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found!")
	}

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	scheduleAndPush(&apiCfg)

	go func() {
		for range ticker.C {
			scheduleAndPush(&apiCfg)
		}
	}()

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func scheduleAndPush(cfg *apiConfig) {
	if os.Getenv("SLACK_API_TOKEN") == "" || os.Getenv("SLACK_CHANNEL") == "" {
		log.Println("Slack cron job skipped, required environment variables missing.")
		return
	}

	fmt.Println()
	log.Println("\033[1;32m=====[ SLACK CRON JOB ]=====\033[0m")

	cfg.scheduleNotifications(context.Background())

	if time.Now().UTC().Hour() == 5 {
		cfg.pushBirthdayNotification(context.Background())
	}

	log.Println("\033[1;32m=====[   CRON EXITED  ]=====\033[0m")
	fmt.Println()
}
