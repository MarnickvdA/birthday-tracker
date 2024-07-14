package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func PostBirthdaySlackMessage(birthdays []Birthday) (string, error) {
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

	return body, nil
}

func PostRequest(token string, body []byte) {
	req, err := http.NewRequest(http.MethodPost, "https://slack.com/api/chat.postMessage", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)
}
