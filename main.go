package main

import (
	"fmt"
	"log"
	"os"
)
import "github.com/joho/godotenv"

func main() {
	app := App{}

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading environment variables")
		log.Fatal(err)
	}

	app.Initialize(os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASS"), os.Getenv("MONGO_BASE_URL"))
	app.Listen("80")
}
