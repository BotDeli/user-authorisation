FROM golang:latest
WORKDIR /app
COPY . .
EXPOSE 50053
CMD go run cmd/main/main.go