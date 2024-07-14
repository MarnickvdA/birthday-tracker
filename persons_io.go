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
	"strings"
	"sync"
	"time"
)

const file = "data/birthdays.csv"

var (
	mu sync.Mutex
)

func initDatabase() {
	mu.Lock()
	defer mu.Unlock()

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		log.Println("No database file found, trying to create new one at " + file)
		err := os.Mkdir("data", fs.ModePerm)

		if err != nil {
			panic(err)
		}

		_, err = os.Create(file)
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

	bs, err := os.ReadFile(file)
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

	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	err = removeLineFromFile(file, file+string(time.Now().Unix()), id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func removeLineFromFile(inputFile, outputFile, id string) error {
	// Open the input file for reading
	inFile, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// Create the output file for writing
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create a scanner to read the input file
	scanner := bufio.NewScanner(inFile)
	writer := bufio.NewWriter(outFile)

	// Read the input file line by line
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line matches the string to be removed
		if strings.Split(strings.TrimSpace(line), ";")[0] != id {
			// Write the line to the output file if it does not match
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return err
	}

	// Flush the writer to ensure all data is written to the output file
	if err := writer.Flush(); err != nil {
		return err
	}

	// Replace the original file with the modified file
	err = os.Rename(outputFile, inputFile)
	if err != nil {
		return err
	}

	return nil
}
