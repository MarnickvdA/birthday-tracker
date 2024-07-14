# Specifies a parent image
FROM golang:1.22.5
 
WORKDIR /app
 
COPY go.mod ./ 
# COPY go.mod go.sum ./
# RUN go mod download

COPY *.go *.html ./
 
RUN CGO_ENABLED=0 GOOS=linux go build -o /birthdays-tracker
 
EXPOSE 1337
 
CMD [ "/birthdays-tracker" ]