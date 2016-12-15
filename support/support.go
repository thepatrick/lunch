package support

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type errorStruct struct {
	Message string `json:"message"`
}

// ErrorWithJSON generates a JSON error response
func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	respBody, err := json.MarshalIndent(errorStruct{message}, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, code)
}

// ResponseWithJSON generates a JSON response (maybe successful)
func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

// Logging Generates simple logging to stdout
func Logging(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request: %v %v\n", r.Method, r.URL)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
