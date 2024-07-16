package main

import (
	"birthdays-tracker/internal/database"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Birthday struct {
	Name string
	Age  int
}

func (cfg *apiConfig) handlerSendBirthdayMessage(w http.ResponseWriter, r *http.Request) {
	persons, err := cfg.DB.ListPersons(r.Context())
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	birthdays := make([]Birthday, 0)
	for _, person := range persons {
		if isBirthday(person) {
			birthdays = append(birthdays, Birthday{Name: person.Name, Age: getAge(person)})
		}
	}

	message, err := sendBirthdaySlackMessage(birthdays)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "%v", message)
}

func sendBirthdaySlackMessage(birthdays []Birthday) (string, error) {
	if len(birthdays) == 0 {
		return "", errors.New("no birthdays to show")
	}

	token, channel := os.Getenv("SLACK_API_TOKEN"), os.Getenv("SLACK_CHANNEL")

	if token == "" || channel == "" {
		log.Println("environment variables SLACK_API_TOKEN and/or SLACK_CHANNEL not found!")
	}

	messages := make([]string, 0)
	for _, b := range birthdays {
		messages = append(messages, fmt.Sprintf(`{
							"type": "rich_text_section",
							"elements": [
								{
									"type": "text",
									"text": "%s turned %d"
								}
							]
						}`, b.Name, b.Age))
	}

	body := fmt.Sprintf(`{
	"channel": "%s",
	"blocks": [
		{
			"type": "header",
			"text": {
				"type": "plain_text",
				"text": "🎉🍾⚡️",
				"emoji": true
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "It's a special day today! We have %d lucky birthday people 🎈"
			}
		},
		{
			"type": "rich_text",
			"elements": [
				{
					"type": "rich_text_list",
					"style": "bullet",
					"indent": 0,
					"border": 0,
					"elements": [
						%s
					]
				}
			]
		}
	]	
}`, channel, len(messages), strings.Join(messages, ",\n"))

	req, err := http.NewRequest(http.MethodPost, "https://slack.com/api/chat.postMessage", bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return "", err
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return "", err
	}

	return string(resBody), nil
}

func getBirthdayParts(p database.Person) (year, month, day int) {
	birthdayParts := strings.Split(p.BirthDate, "-")

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

func isBirthday(p database.Person) bool {
	_, m, d := getBirthdayParts(p)

	return int(time.Now().Month()) == m && time.Now().Day() == d
}

func getAge(p database.Person) int {
	y, m, d := getBirthdayParts(p)

	age := time.Now().Year() - y

	if time.Now().Compare(time.Date(time.Now().Year(), time.Month(m), d, 0, 0, 0, 0, time.Local)) < 0 {
		age--
	}

	return age
}
