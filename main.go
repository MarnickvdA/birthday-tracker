package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Person struct {
	ID        int
	Name      string
	Birthdate string
}

type PersonPageData struct {
	Persons []Person
}

var (
	persons []Person
	counter = 0
	tmpl    = template.Must(template.ParseFiles("index.html"))
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	data := PersonPageData{
		Persons: persons,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func postBirthdayHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	name, birthdate := r.Form.Get("name"), r.Form.Get("birthdate")

	if name != "" && birthdate != "" {
		addPerson(name, birthdate)
	}

	getHandler(w, r)
}

func removeBirthdayHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	removePerson(id)

	fmt.Fprintln(w, "OK")
}

func addPerson(name, birthdate string) {
	persons = append(persons, Person{ID: counter, Name: name, Birthdate: birthdate})
	counter++
}

func removePerson(id int) {
	for i, person := range persons {
		if person.ID == id {
			persons = append(persons[:i], persons[i+1:]...)
		}
	}

	fmt.Println(persons)
}

func main() {
	addPerson("Marnick", "1998-04-16")
	addPerson("Simon", "1995-09-10")

	http.HandleFunc("GET /", getHandler)
	http.HandleFunc("POST /", postBirthdayHandler)
	http.HandleFunc("DELETE /{id}", removeBirthdayHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
