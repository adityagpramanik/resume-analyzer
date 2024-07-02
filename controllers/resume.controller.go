package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const MAX_UPLOAD_SIZE = int64(10 * 1024 * 1024) // 10 MB limit

// region UploadHandler
func UploadHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
  }

  err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)
  if err != nil {
    log.Println("Error parsing multipart form:", err)
    http.Error(w, "Error parsing upload", http.StatusBadRequest)
    return
  }

  file, _, err := r.FormFile("file")
  if err != nil {
    log.Println("Error getting uploaded file:", err)
    http.Error(w, "Error getting uploaded file", http.StatusBadRequest)
    return
  }

  defer file.Close()

  data, err := io.ReadAll(file)
  if err != nil {
    log.Println("Error reading uploaded file:", err)
    http.Error(w, "Error reading uploaded file", http.StatusInternalServerError)
    return
  }

  // http.Post("http://localhost:8000/extracter/api/v1/extractResume",)

  processedData := string(data)
  fmt.Fprintf(w, "Successfully uploaded and processed file: %s", processedData)
}