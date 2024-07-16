FROM golang:alpine AS build
 
WORKDIR /app

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
 
# COPY go.mod ./ 
COPY . .
RUN go mod download
 
RUN CGO_ENABLED=0 GOOS=linux go build -o /birthday-tracker

FROM scratch

COPY --from=build /birthday-tracker .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# copy all other files
COPY .env.prod .env
COPY ./web ./web
COPY ./static ./static

CMD ["/birthday-tracker"]