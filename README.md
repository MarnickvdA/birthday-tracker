# Birthday Tracker 🎈

Web app written in Go to track birthdays of people.

**Features**:

- See overview of persons and their birthdays
- Add new person with birthday
- Remove persons
- (WIP) Trigger Slack message of birthdays happening today
- (WIP) Add cron configuration for running the birthday check every day

**Technical improvements**:

- (WIP) Add automatic db migrations to the docker-compose script w/ goose

## Dependencies

- docker / docker-compose
- go
- goose
- sqlc

## Setup

We have some environment variables on which we depend in this project. Execute the command below and update the variables.

```bash
cp .env.example .env
```

Most important is to have a PostgresQL server running to which we can connect. Be sure to add the connection URL as `DB_URL` in .env

```bash
# Run goose migrations on the postgres db
goose postgres $DB_URL up
```

## Running with Docker

```bash
docker compose up -d --build
```
