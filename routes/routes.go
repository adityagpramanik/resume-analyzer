package routes

import (
	"github.com/gorilla/mux"
	"resume-scanner.com/resume-scanner/controllers"
)

func SetupRoutes() *mux.Router {
    router := mux.NewRouter()

    // MARK: Root routes
    router.HandleFunc("/", controllers.RootHandler).Methods("GET")

    // MARK: Resume routes
    router.HandleFunc("/resumes/analyze", controllers.ResumeAnaylze).Methods("POST")
    router.HandleFunc("/resumes/analyze", controllers.ResumeAnalyzeCloudFile).Methods("GET")
    router.HandleFunc("/resumes/upload", controllers.ResumeUpload).Methods("POST")

    return router
}
