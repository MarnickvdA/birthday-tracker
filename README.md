# Birthday Tracker 🎈

Birthday tracker integrated with Slack API to notify a channel about people's birthdays! It's a Web app written in Go to track birthdays of people. 

**Features**:

- See overview of persons and their birthdays
- Add new person with birthday
- Remove persons
- Automatically pushing Slack message of today's birthdays
- (TODO) Edit person

**Technical Features**:

- Web Server running fully on Go's `net/http` standard library
- docker-compose for automatic container deployments
- Database migrations ran on app startup
- Frontend templated with Go's `html/template` standard library
- Minimal CSS for builds with TailwindCSS

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
