# Birthday Tracker 🎈

Birthday tracker integrated with Slack API to notify a channel about people's birthdays! It's a Web app written in Go to track birthdays of people.

**Features**:

- [x] See overview of persons and their birthdays
- [x] Add new person with birthday
- [x] Remove persons
- [x] Automatically pushing Slack message of today's birthdays
- [ ] Edit person
- [ ] Automatical birthday messages with LLM generated text
- [ ] Improved overview of birthdays in calendar view

**Technical Features**:

- [x] Web Server running fully on Go's `net/http` standard library
- [x] docker-compose for automatic container deployments
- [x] Database migrations ran on app startup
- [x] Frontend templated with Go's `html/template` standard library
- [x] Minimal CSS for builds with TailwindCSS
- [ ] Setup CI/CD Pipeline on GitHub to build image
- [ ] Use HTMX because why not?
- [ ] Add telemetry setup

## Dependencies

- docker / docker-compose
- go
- goose
- sqlc
- npm / tailwindcss

## Setup

We have some environment variables on which we depend in this project. Execute the command below and update the variables.

```bash
cp .env.example .env
go mod download
```

Most important is to have a PostgresQL server running to which we can connect. Be sure to add the connection URL as `DB_URL` in .env and .env.prod. On app startup, migrations from the `sql/schema` folder will automatically be run.

## Local development

For streamlined development, always use our special dev script. It will load all the necessary resources for running the dev environment.

```bash
./dev.sh
```

## Running with Docker

```bash
docker compose up -d --build
```
