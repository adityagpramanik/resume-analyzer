package controllers

import (
	"fmt"
	"net/http"
)

func RootHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(res, "Hello world!\n")
}