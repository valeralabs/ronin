package events

import (
	"github.com/gorilla/mux"
	"github.com/syvita/ronin/events/mempool"
	"github.com/syvita/ronin/shared"

	"encoding/json"
	"github.com/syvita/ronin/db"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var Database *db.Database

var logger = log.New(os.Stderr, "[EVENT]: ", log.Ldate|log.Ltime|log.Lshortfile)

func init() {
	mempool.Logger = logger
}

func NewMempoolTX(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		writer.WriteHeader(500)
		shared.Send(writer, shared.Object{"error": "failed to read body"})
		return
	}

	if len(body) == 0 {
		writer.WriteHeader(400)
		shared.Send(writer, shared.Object{"error": "body is required"})
		return
	}

	var TXs []string

	err = json.Unmarshal(body, &TXs)

	if err != nil {
		writer.WriteHeader(400)
		shared.Send(writer, shared.Object{"error": "json is invalid"})
		return
	}

	TXIDs, err := mempool.NewTX(TXs)

	if err != nil {
		writer.WriteHeader(500)
		shared.Send(writer, shared.Object{"error": err.Error()})
		return
	}

	shared.Send(writer, shared.Object{"ok": true, "txids": TXIDs})
}

func Listen(address string) {
	logger.Println("starting...")

	router := mux.NewRouter()

	router.HandleFunc("/new_mempool_tx", NewMempoolTX).Methods("POST")

	err := http.ListenAndServe(address, router)

	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
}
