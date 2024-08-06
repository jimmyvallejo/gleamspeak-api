package common

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ContextKey string

const UserContextKey ContextKey = "user"