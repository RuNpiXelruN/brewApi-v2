package utils

// Result type
type Result struct {
	Success *Success `json:"success"`
	Error   *Error   `json:"error"`
}

// Success type
type Success struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// Error type
type Error struct {
	Status     int    `json:"status"`
	StatusText string `json:"status_text"`
}
