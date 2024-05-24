package util

type ValidationError struct {
	Message  string `json:"mensage"`
	Field    string `json:"field"`
	Location string `json:"local"`
	Help     string `json:"help"`
}
