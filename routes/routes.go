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
    router.HandleFunc("/upload", controllers.UploadHandler).Methods("POST")

    return router
}
