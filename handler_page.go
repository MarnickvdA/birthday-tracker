package main

import (
	"birthday-tracker/internal/database"
	"html/template"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerGetHomePage(w http.ResponseWriter, r *http.Request) {
	persons, err := cfg.DB.ListPersons(r.Context())
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	type PersonPageData struct {
		Persons []database.Person
	}

	tmpl, err := template.ParseFiles("./web/homepage.html")

	if err != nil {
		log.Fatal(err)
	}

	if err := tmpl.Execute(w, PersonPageData{
		Persons: persons,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
