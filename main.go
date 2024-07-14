package main

import (
	"birthdays-tracker/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

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

	http.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css"))))

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found!")
	}

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
