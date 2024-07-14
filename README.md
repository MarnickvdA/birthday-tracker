# Basic Birthday tracker

Super simple Go app to track birthdays of people.

## Running in Docker

```bash
docker build --rm -t birthdays-tracker:alpha .
docker run -d -p 1337:1337 --name birthdays-tracker birthdays-tracker:alpha 
```

## TODO

- Add Postgres DB
- Create docker-compose workflow to enable launching the app with the PostgresDB
- Add endpoint to check who's birthday it is TODAY!
- Add cron job integration for Slack notifications
- Add cron job integration with reminder emails
