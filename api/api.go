package api

import (
	"github.com/gorilla/mux"
	"github.com/syvita/ronin/db"
	"github.com/syvita/ronin/shared"

	"log"
	"net/http"
	"os"
)

var Database *db.Database

var logger = log.New(os.Stderr, "[API]: ", log.Ldate|log.Ltime|log.Lshortfile)

func status(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	shared.Send(writer, shared.Object{"ok": true, "version": "TODO"})
}

func Listen(address string) {
	logger.Println("starting...")

	router := mux.NewRouter()

	router.HandleFunc("/status", status).Methods("GET")

	err := http.ListenAndServe(address, router)

	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
}
