#!/bin/bash
git pull origin/main

docker build --rm -t birthdays-tracker:alpha .

docker container stop birthdays-tracker
docker container rm birthdays-tracker

docker run -d -p 1337:1337 --name birthdays-tracker birthdays-tracker:alpha

curl -f -s -o /dev/null -v 'http://localhost:1337'
 