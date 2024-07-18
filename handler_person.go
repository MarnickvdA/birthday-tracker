package main

import (
	"birthday-tracker/internal/database"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreatePerson(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	name, birthdate := r.Form.Get("name"), r.Form.Get("birthdate")
	if name == "" || birthdate == "" {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	person, err := cfg.DB.CreatePerson(r.Context(), database.CreatePersonParams{
		Name:      name,
		BirthDate: birthdate,
	})

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Couldn't save person", http.StatusInternalServerError)
		return
	}

	log.Printf("Added %v\n", person)

	cfg.handlerGetHomePage(w, r)
}

func (cfg *apiConfig) handlerRemovePerson(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := cfg.DB.DeletePerson(r.Context(), id)

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Could not remove person", http.StatusInternalServerError)
		return
	}

	log.Printf("Removed person with ID: %v\n", id)
}
