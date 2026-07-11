FROM golang:1.26-alpine

RUN apk add --no-cache make

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum .air.toml ./
RUN go mod download

COPY . .

EXPOSE 6379 40000

CMD ["air"]
