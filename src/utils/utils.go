package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

var result Result

func ParseTime(timeRaw string) time.Time {
	const timeLayout = "02-01-2006 15:04 (MST)"
	t, _ := time.Parse(timeLayout, timeRaw)
	return t
}

// for handling Responses in middleware only
func Response(w http.ResponseWriter, result *Result) {
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
