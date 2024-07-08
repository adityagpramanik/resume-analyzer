package main

// MARK: all imports goes here
import (
	"fmt"
	"log"
	"net/http"

	"resume-scanner.com/resume-scanner/routes"
)

func main() {
	router := routes.SetupRoutes()
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
