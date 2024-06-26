package main

import (
	"RealTimeForum/Server"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	portEnv := os.Getenv("PORT")
	if portEnv == "" {
		portEnv = "8080"
	}
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		log.Printf("Error getting port: %s\n", err.Error())
		os.Exit(1)
	}
	// Ports can only be > 0 && <= 65535
	if !(port > 0 && port <= math.MaxUint16) {
		log.Printf("Invalid port: %s\n", portEnv)
		os.Exit(1)
	}
	Server.GoLive(portEnv)
}
