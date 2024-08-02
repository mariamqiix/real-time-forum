# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o sandbox-bin main.go

EXPOSE 8080
ENV PORT=8080

CMD ["./sandbox-bin"]