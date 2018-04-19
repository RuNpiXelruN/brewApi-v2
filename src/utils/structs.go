package utils

type Result struct {
	Success *Success `json:"success"`
	Error   *Error   `json:"error"`
}

type Success struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}

type Error struct {
	Status int `json:"status"`
}
