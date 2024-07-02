package main

//all imports goes here
import (
	"fmt"
	"log"
	"net/http"

	resumecontroller "resume-scanner.com/resume-scanner/controllers"
)

func rootHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(res, "Hello world!\n")
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/upload", resumecontroller.UploadHandler)
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
