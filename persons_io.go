package main

import (
	"bufio"
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Person struct {
	ID        string
	Name      string
	Birthdate string
}

type Birthday struct {
	Name string
	Age  int
}

var (
	mu sync.Mutex
)

func getDatabaseFile() string {
	path := os.Getenv("DATABASE_PATH")

	if path != "" {
		return path
	} else {
		return "data/birthdays.csv"
	}
}

func initDatabase() {
	mu.Lock()
	defer mu.Unlock()

	if _, err := os.Stat(getDatabaseFile()); errors.Is(err, os.ErrNotExist) {
		log.Println("No database file found, trying to create new one at " + getDatabaseFile())
		err := os.Mkdir("data", fs.ModePerm)

		if err != nil {
			panic(err)
		}

		_, err = os.Create(getDatabaseFile())
		if err != nil {
			panic(err)
		}

		log.Println("Database file created successfully.")
	} else {
		log.Println("Database file detected, we'll use that.")
	}
}

func getPersons() (persons []Person, err error) {
	mu.Lock()
	defer mu.Unlock()

	bs, err := os.ReadFile(getDatabaseFile())
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(bs), "\n")
	persons = make([]Person, 0)
	for _, line := range lines {
		parts := strings.Split(line, ";")

		if len(parts) != 3 {
			continue
		}

		persons = append(persons, Person{parts[0], parts[1], parts[2]})
	}

	return persons, nil
}

func addPerson(name, birthdate string) (person Person, err error) {
	if name == "" || birthdate == "" {
		fmt.Println("Invalid name or birthdate")
		return Person{}, errors.New("invalid name or birthdate")
	}

	hasher := crypto.SHA256.New()
	hasher.Write([]byte(name + birthdate + time.UnixDate))
	id := hex.EncodeToString(hasher.Sum(nil))[:16]

	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile(getDatabaseFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return Person{}, err
	}

	_, err = f.Write([]byte(fmt.Sprintf("%v;%v;%v\n", id, name, birthdate)))
	if err != nil {
		return Person{}, err
	}

	if err := f.Close(); err != nil {
		return Person{}, err
	}

	return Person{id, name, birthdate}, nil
}

func removePerson(id string) (d string, err error) {
	mu.Lock()
	defer mu.Unlock()

	inputFile := getDatabaseFile()
	outputFile := getDatabaseFile() + string(time.Now().Unix())

	inFile, err := os.Open(inputFile)
	if err != nil {
		return "", err
	}
	defer inFile.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(inFile)
	writer := bufio.NewWriter(outFile)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Split(strings.TrimSpace(line), ";")[0] != id {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				return "", err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if err := writer.Flush(); err != nil {
		return "", err
	}

	err = os.Rename(outputFile, inputFile)
	if err != nil {
		return "", err
	}

	return id, nil
}

func getBirthdaysToday() (birthdays []Birthday, err error) {
	persons, err := getPersons()

	if err != nil {
		return nil, err
	}

	birthdays = make([]Birthday, 0)

	for _, person := range persons {
		if person.isBirthday() {
			birthdays = append(birthdays, Birthday{Name: person.Name, Age: person.getAge()})
		}
	}

	return birthdays, nil
}

func (p Person) getBirthdayParts() (year, month, day int) {
	birthdayParts := strings.Split(p.Birthdate, "-")

	y, err := strconv.Atoi(birthdayParts[0])

	if err != nil {
		panic(err)
	}

	m, err := strconv.Atoi(birthdayParts[1])

	if err != nil {
		panic(err)
	}

	d, err := strconv.Atoi(birthdayParts[2])

	if err != nil {
		panic(err)
	}

	return y, m, d
}

func (p Person) isBirthday() bool {
	_, m, d := p.getBirthdayParts()

	return int(time.Now().Month()) == m && time.Now().Day() == d
}

func (p Person) getAge() int {
	y, m, d := p.getBirthdayParts()

	age := time.Now().Year() - y

	if time.Now().Compare(time.Date(time.Now().Year(), time.Month(m), d, 0, 0, 0, 0, time.Local)) < 0 {
		age--
	}

	return age
}
