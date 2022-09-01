package main

import (
	"fmt"
	"log"

	"github.com/fxmbx/go-paystack-pr/cmd/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	server := api.Server{}
	fmt.Println("serving on port 8080")
	err = server.Router()
	if err != nil {
		log.Fatal("Error Serving up router:", err.Error())
	}

}
