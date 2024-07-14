package main

import (
	"log"
	"net/http"
	"text/template"
)

type Person struct {
	ID        string
	Name      string
	Birthdate string
}

type PersonPageData struct {
	Persons []Person
}

func getPageHandler(w http.ResponseWriter, r *http.Request) {
	persons, err := getPersons()

	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := PersonPageData{
		Persons: persons,
	}

	if err := template.Must(template.ParseFiles("index.html")).Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addBirthdayHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	person, err := addPerson(r.Form.Get("name"), r.Form.Get("birthdate"))

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid Input", http.StatusInternalServerError)
	} else {
		log.Printf("Added %v\n", person)
		getPageHandler(w, r)
	}
}

func removeBirthdayHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	d, err := removePerson(id)

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Could not remove birthday", http.StatusInternalServerError)
	} else {
		log.Printf("Removed person with ID: %v\n", d)
	}
}

func main() {
	http.HandleFunc("GET /", getPageHandler)
	http.HandleFunc("POST /", addBirthdayHandler)
	http.HandleFunc("DELETE /{id}", removeBirthdayHandler)

	log.Println("Starting server on port 1337")
	log.Fatal(http.ListenAndServe(":1337", nil))
}
