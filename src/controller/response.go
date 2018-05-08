package controller

import (
	"encoding/json"
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"net/http"
)

// Response func
func Response(w http.ResponseWriter, result *utils.Result) {
	if result.Error != nil {
		w.WriteHeader(result.Error.Status)
		w.Header().Set("Status", http.StatusText(result.Error.Status))

		data, _ := json.Marshal(result.Error)
		w.Write(data)
		return
	}

	w.WriteHeader(result.Success.Status)
	w.Header().Set("Status", http.StatusText(result.Success.Status))

	data, _ := json.Marshal(result.Success.Data)
	w.Write(data)
}
