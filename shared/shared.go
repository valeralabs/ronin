package shared

import (
	"encoding/json"
	"net/http"
)

type Object map[string]interface{}

func Send(writer http.ResponseWriter, body Object) {
	err := json.NewEncoder(writer).Encode(body)

	if err != nil {
		writer.Write([]byte("failed to serialize error"))
	}
}
