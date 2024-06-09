# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project to the working directory
COPY go.mod ./go.mod
COPY cmd ./cmd
COPY internal ./internal
COPY www ./www
COPY forum-db.sqlite ./forum-db.sqlite

RUN go mod tidy
# Build the Go application
RUN go build -o sandbox-bin cmd/sandbox/main.go

EXPOSE 8080
ENV PORT=8080

# Command to run the executable
CMD ["./sandbox-bin"]
