FROM golang:alpine AS build
 
WORKDIR /app

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
 
RUN apk add --update npm

# COPY go.mod ./ 
COPY . .
RUN go mod download
 
RUN npm ci --prefix web/tailwindcss
RUN npm run --prefix web/tailwindcss build

RUN CGO_ENABLED=0 GOOS=linux go build -o /birthday-tracker

FROM scratch

COPY --from=build /birthday-tracker .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/static/ /static/

# copy all other files
COPY .env.prod .env
COPY ./web ./web

CMD ["/birthday-tracker"]