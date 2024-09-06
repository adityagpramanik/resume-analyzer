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
	commonservices "resume-scanner.com/resume-scanner/services/common"
	"resume-scanner.com/resume-scanner/services/files"
)

// Local models here
type CloudAnalyserScrapperReqDTO struct {
	Url string `json:"url"`
}

const MAX_UPLOAD_SIZE = int64(10 * 1024 * 1024) // 10 MB limit

func init() {
	godotenv.Load()
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
	parserReq, err := http.NewRequest("POST", RESUME_HELPER_SERVICE+"extracter/api/v1/extractResume", &requestBody)
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
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// read params
	params := req.URL.Query()

	file_id, present := params["file_id"]
	if !present || len(file_id) == 0 {
		http.Error(res, "file_id not present", http.StatusBadRequest)
	}

	// get file info from the db
	file, err := files.GetFileByFileId(file_id[0])
	if err != nil {
		log.Fatal("error: ", err)
		http.Error(res, "Unable to fetch file.", http.StatusBadRequest)
	}

	// get public url from the file from db
	RESUME_BUCKET := os.Getenv("RESUME_BUCKET")
	url, err := commonservices.GetFileUrl(RESUME_BUCKET, file.FileName)
	if err != nil {
		log.Println("Error public url for the document:", err)
		http.Error(res, "Error analyzing file from the url.", http.StatusBadRequest)
		return
	}
	// request scrapper-server to analyze the file
	RESUME_HELPER_SERVICE := os.Getenv("RESUME_HELPER_SERVICE")
	requestBody := CloudAnalyserScrapperReqDTO{
		Url: url,
	}

	// Buffer to hold the JSON data
	var bodyBuff bytes.Buffer

	// Create a new JSON encoder and disable HTML escaping
	encoder := json.NewEncoder(&bodyBuff)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	parserReq, err := http.NewRequest("POST", RESUME_HELPER_SERVICE+"extracter/api/v1/extractResume", &bodyBuff)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	// Set headers
	parserReq.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(parserReq)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		http.Error(res, "Unable to process file: "+file_id[0], resp.StatusCode)
		return
	}
	fmt.Fprint(res, string(respBody))
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

	if handler.Header["Content-Type"][0] != "application/pdf" {
		http.Error(res, "Error getting uploaded file", http.StatusBadRequest)
		return
	}

	RESUME_BUCKET := os.Getenv("RESUME_BUCKET")
	filename := uuid.New().String() + ".pdf";

	err = commonservices.UploadFileBuffer(RESUME_BUCKET, filename, file, handler.Size);
	if err != nil {
		log.Println("Error getting uploaded file:", err)
		http.Error(res, "Error getting uploaded file", http.StatusBadRequest)
		return
	}

	url, err := commonservices.GetFileUrl(RESUME_BUCKET, filename)
	if err != nil {
		http.Error(res, "Unable to generate file link", http.StatusBadRequest)
	}

	session := commonservices.GetSession()

	err = session.Query("insert into resume_analyzer.files (file_id, file_name, bucket, url, createdat, updatedat) VALUES (uuid(), ?, ?, ?, dateof(now()), dateof(now()));", filename, RESUME_BUCKET, url).Exec()

	if err != nil {
		http.Error(res, "Error processing file.", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"url": url,
	}

	result, err := json.Marshal(data)
	if err != nil {
		http.Error(res, "Unable to parse data", http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(result)
	if err != nil {
		http.Error(res, "Error getting uploaded resume", http.StatusBadRequest)
	}
}
