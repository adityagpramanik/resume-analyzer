package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	commonservices "resume-scanner.com/resume-scanner/common.services"
)

const MAX_UPLOAD_SIZE = int64(10 * 1024 * 1024) // 10 MB limit

func init() {
	godotenv.Load();
}

// MARK: Analyze attached file
func ResumeAnaylze(res http.ResponseWriter, req *http.Request) {
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

	RESUME_HELPER_SERVICE := os.Getenv("RESUME_HELPER_SERVICE")
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

// MARK: Analyze file from blob server
func ResumeAnalyzeCloudFile(res http.ResponseWriter, req *http.Request) {
	// TODO: implement this
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(res, "Analyzed resume")
}

// MARK: Upload resume to blob storage server
func ResumeUpload(res http.ResponseWriter, req *http.Request) {
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

	file, handler, err := req.FormFile("file")
	if err != nil {
		log.Println("Error getting uploaded file:", err)
		http.Error(res, "Error getting uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	RESUME_BUCKET := os.Getenv("RESUME_BUCKET");
	filename := uuid.New()

	err = commonservices.UploadFileBuffer(RESUME_BUCKET, filename.String(), file, handler.Size);
	if err != nil {
		log.Println("Error getting uploaded file:", err)
		http.Error(res, "Error getting uploaded file", http.StatusBadRequest)
		return
	}

	url, err := commonservices.GetFileUrl(RESUME_BUCKET, filename.String())
	if err != nil {
		http.Error(res, "Unable to generate file link", http.StatusBadRequest)
	}

	data := map[string]interface{}{
		"url": url,
	}

	result, err := json.Marshal(data)
	if err != nil {
		http.Error(res, "Unable to parse data", http.StatusInternalServerError)
	}
	
	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(result);
	if err != nil {
		http.Error(res, "Error getting uploaded resume", http.StatusBadRequest)
	}
}