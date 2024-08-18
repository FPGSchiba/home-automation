package util

type Response struct {
	Message string                 `json:"message"`
	Status  string                 `json:"status"`
	Data    map[string]interface{} `json:"data"`
}
