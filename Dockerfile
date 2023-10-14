# The base go-image
FROM golang:1.21.3-bullseye

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o /check-websites

EXPOSE 8080

CMD [ "/check-websites" ]