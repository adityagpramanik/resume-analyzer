package main

// MARK: all imports goes here
import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"resume-scanner.com/resume-scanner/routes"
)

func init() {
	err := godotenv.Load();
	if err != nil {
		log.Fatal("Error loading environment secrets");
	}
}

func main() {
	router := routes.SetupRoutes()
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
