package common

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Error: message})
}

// WriteProblem writes an RFC 7807 Problem Detail response
func WriteProblem(w http.ResponseWriter, problem interface{}) {
	w.Header().Set("Content-Type", "application/problem+json")

	// Extract status from problem
	status := 500
	if p, ok := problem.(*ProblemDetail); ok {
		status = p.Status
	} else if p, ok := problem.(*ValidationProblem); ok {
		status = p.Status
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(problem)
}
