package main

import (
	"birthdays-tracker/internal/database"
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

//go:embed static
var static embed.FS

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

	http.HandleFunc("GET /", apiCfg.handlerGetHomePage)
	http.HandleFunc("POST /", apiCfg.handlerCreatePerson)
	http.HandleFunc("DELETE /{id}", apiCfg.handlerRemovePerson)
	http.HandleFunc("POST /today", apiCfg.handlerSendBirthdayMessage)

	//go:generate npm run build
	http.Handle("GET /static/", http.FileServer(http.FS(static)))

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found!")
	}

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			log.Println("What is that mysterious ticking noise???")
		}
	}()

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
