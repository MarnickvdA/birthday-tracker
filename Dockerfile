FROM golang:1.22.5 AS build
 
WORKDIR /app
 
# COPY go.mod ./ 
COPY . .
RUN go mod download
 
RUN CGO_ENABLED=0 GOOS=linux go build -o /birthdays-tracker

FROM scratch

COPY --from=build /birthdays-tracker .

# copy all other files
COPY .env.prod .env
COPY ./web ./web

CMD ["/birthdays-tracker"]