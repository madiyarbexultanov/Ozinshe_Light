FROM golang:1.23.3-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .env .
COPY . .

RUN go build -o /ozinshe-go

EXPOSE 8080

ENTRYPOINT [ "/ozinshe-go" ]
