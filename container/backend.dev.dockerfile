FROM golang:1.25

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY src/backend/go.mod src/backend/go.sum ./
RUN go mod download

WORKDIR /app

CMD ["air", "-c", ".air.toml"]
