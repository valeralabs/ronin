package main

import (
	"github.com/gorilla/mux"
	"github.com/syvita/ronin/db"

	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ClientHandler struct {
	Database *db.Database
	Address  string
}

type EventHandler struct {
	Database *db.Database
	Address  string
}

type Object map[string]interface{}

var RedisAddr = "localhost:6379"

var logger = log.New(os.Stderr, "[MAIN]: ", log.Ldate|log.Ltime|log.Lshortfile)

func send(writer http.ResponseWriter, body Object) {
	err := json.NewEncoder(writer).Encode(body)

	if err != nil {
		writer.Write([]byte("failed to serialize error"))
	}
}

//TODO: move this to a dedicated module
func (handler ClientHandler) Start() {
	logger := log.New(os.Stderr, "[CLIENT HANDLER]: ", log.Ldate|log.Ltime|log.Lshortfile)
	router := mux.NewRouter()

	router.HandleFunc("/points/{username}", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")

		username := mux.Vars(request)["username"]

		logger.Printf("finding user \"%s\"\n", username)

		if username == "" {
			send(writer, Object{"error": "username is required"})
			return
		}

		user, err := handler.Database.GetUser(username)

		if err != nil {
			if err == db.ErrNil {
				writer.WriteHeader(404)
				send(writer, Object{"error": "No record found for " + username})
				return
			}

			writer.WriteHeader(500)
			send(writer, Object{"error": "No record"})
			return
		}

		writer.WriteHeader(200)
		send(writer, Object{"user": user})
	}).Methods("GET")

	router.HandleFunc("/points", func(writer http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)

		if err != nil {
			writer.WriteHeader(500)
			send(writer, Object{"error": "failed to read body"})
			return
		}

		if len(body) == 0 {
			writer.WriteHeader(400)
			send(writer, Object{"error": "body is required"})
			return
		}

		var user db.User

		json.Unmarshal(body, &user)

		err = handler.Database.SaveUser(&user)

		if err != nil {
			writer.WriteHeader(500)
			send(writer, Object{"error": "failed to save user"})
			return
		}

		logger.Println("storing user", user.Username)

		send(writer, Object{"user": user})
	}).Methods("POST")

	err := http.ListenAndServe(handler.Address, router)

	if err != nil {
		panic(err)
	}
}

//TODO: move this to a dedicated module
func (handler EventHandler) Start() {
	logger := log.New(os.Stderr, "[EVENT HANDLER]: ", log.Ldate|log.Ltime|log.Lshortfile)
	router := mux.NewRouter()

	router.HandleFunc("/points/{username}", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")

		username := mux.Vars(request)["username"]

		logger.Printf("finding user \"%s\"\n", username)

		if username == "" {
			send(writer, Object{"error": "username is required"})
			return
		}

		user, err := handler.Database.GetUser(username)

		if err != nil {
			if err == db.ErrNil {
				writer.WriteHeader(404)
				send(writer, Object{"error": "No record found for " + username})
				return
			}

			writer.WriteHeader(500)
			send(writer, Object{"error": "No record"})
			return
		}

		writer.WriteHeader(200)
		send(writer, Object{"user": user})
	}).Methods("GET")

	router.HandleFunc("/points", func(writer http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)

		if err != nil {
			writer.WriteHeader(500)
			send(writer, Object{"error": "failed to read body"})
			return
		}

		if len(body) == 0 {
			writer.WriteHeader(400)
			send(writer, Object{"error": "body is required"})
			return
		}

		var user db.User

		json.Unmarshal(body, &user)

		err = handler.Database.SaveUser(&user)

		if err != nil {
			writer.WriteHeader(500)
			send(writer, Object{"error": "failed to save user"})
			return
		}

		logger.Println("storing user", user.Username)

		send(writer, Object{"user": user})
	}).Methods("POST")

	err := http.ListenAndServe(handler.Address, router)

	if err != nil {
		panic(err)
	}
}

func main() {
	database, err := db.NewDatabase(RedisAddr)

	if err != nil {
		logger.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	logger.Println("Connected to Redis successfully")

	client := ClientHandler{database, ":3999"}
	event := EventHandler{database, ":3700"}

	// goroutines >>>>> anything else

	go client.Start()
	go event.Start()

	for {
	}
}
