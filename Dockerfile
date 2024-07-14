FROM golang:1.22.5 AS build
 
WORKDIR /app
 
# COPY go.mod ./ 
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
 
RUN CGO_ENABLED=0 GOOS=linux go build -o /birthdays-tracker

FROM scratch

COPY --from=build /birthdays-tracker .

# copy all other files
COPY .env .
COPY /templates /templates

CMD ["/birthdays-tracker"]