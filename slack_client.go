package main

import (
	"birthday-tracker/internal/database"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func sendBirthdaySlackMessage(bd []database.GetScheduledBirthdayNotificationsForTodayRow) (string, error) {
	if len(bd) == 0 {
		return "", errors.New("no birthdays to show")
	}

	token := os.Getenv("SLACK_API_TOKEN")
	if token == "" {
		return "", errors.New("environment variable SLACK_API_TOKEN not found")
	}

	body, err := composeSlackMessage(bd)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "https://slack.com/api/chat.postMessage", bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Printf("client: could not create request: %s\n", err)
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("client: error making http request: %s\n", err)
		return "", err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("client: could not read response body: %s\n", err)
		return "", err
	}

	return string(resBody), nil
}

func composeSlackMessage(birthdays []database.GetScheduledBirthdayNotificationsForTodayRow) (string, error) {
	channel := os.Getenv("SLACK_CHANNEL")
	if channel == "" {
		return "", errors.New("environment variable SLACK_CHANNEL not found")
	}

	messages := make([]string, 0)
	for _, b := range birthdays {
		messages = append(messages, fmt.Sprintf(`{
							"type": "rich_text_section",
							"elements": [
								{
									"type": "text",
									"text": "%s turned %s"
								}
							]
						}`, b.Name, b.Age))
	}

	return fmt.Sprintf(`{
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
}`, channel, len(messages), strings.Join(messages, ",\n")), nil
}
