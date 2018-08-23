package utils

import (
	"encoding/json"
	"net/http"
)

// Respond func
// for handling Responses in middleware only
func Respond(w http.ResponseWriter, result *Result) {
	if result.Error != nil {
		w.WriteHeader(result.Error.StatusCode)
		w.Header().Set("Status", http.StatusText(result.Error.StatusCode))

		data, _ := json.Marshal(result.Error)
		w.Write(data)
		return
	}

	w.WriteHeader(result.Success.StatusCode)
	w.Header().Set("Status", http.StatusText(result.Success.StatusCode))

	data, _ := json.Marshal(result.Success.Data)
	w.Write(data)
}
