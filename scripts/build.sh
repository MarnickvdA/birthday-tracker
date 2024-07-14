#!/bin/bash
git checkout main
git pull

docker build . --rm -t birthdays-tracker:alpha -f ./docker/Dockerfile

docker container stop birthdays-tracker
docker container rm birthdays-tracker

docker run -d -p 1337:1337 --name birthdays-tracker birthdays-tracker:alpha

curl -f -s -o /dev/null -v 'http://localhost:1337'
 