package models

// ErrorResponse is a struct that contains the error response to be sent to the client.
// It contains the status code, title, and detail, as well as the timestamp and path of the request.
type ErrorResponse struct {
	Status    int    `json:"status"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	Path      string `json:"path"`
	Timestamp int64  `json:"timestamp"`
}
