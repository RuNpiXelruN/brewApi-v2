package utils

type Result struct {
	Success *Success `json:"success"`
	Error   *Error   `json:"error"`
}

type Success struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type Error struct {
	Status     int    `json:"status"`
	StatusText string `json:"status_text"`
}
