#!/bin/bash

# Do a fresh install if tailwindcss hasn't been initialized.
if [ ! -d "web/tailwindcss/node_modules" ]; then
    echo -e "==>\033[0;33m Required node_modules dependencies not found, we'll install that now...\033[0m"
    npm install --prefix web/tailwindcss
    echo ""
fi

echo -e "==>\033[1;32m Compiling tailwindcss based on styling in web/*.html\033[0m"
npm run --prefix web/tailwindcss build

echo -e "\n==>\033[1;32m Starting the Go application\033[0m"
go run *.go