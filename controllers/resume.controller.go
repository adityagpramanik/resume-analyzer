package controllers

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

const RESUME_HELPER_SERVICE = "http://localhost:8000/extracter/"
const MAX_UPLOAD_SIZE = int64(10 * 1024 * 1024) // 10 MB limit

// MARK: UploadHandler
func UploadHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseMultipartForm(MAX_UPLOAD_SIZE)
	if err != nil {
		log.Println("Error parsing multipart form:", err)
		http.Error(res, "Error parsing upload", http.StatusBadRequest)
		return
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		log.Println("Error getting uploaded file:", err)
		http.Error(res, "Error getting uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	writer.Close()

	parserReq, err := http.NewRequest("POST", RESUME_HELPER_SERVICE+"api/v1/extractResume", &requestBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	parserReq.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(parserReq)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(res, "Unable to parse resume", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, err = io.Copy(res, resp.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
