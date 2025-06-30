package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// HTTPError represents an HTTP error
type HTTPError struct {
	Message string
	Status  int
}

func (e *HTTPError) Error() string {
	return e.Message
}

// DecodeJSONBody decodes the JSON body of a request into the given interface
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return &HTTPError{
			Message: "Content-Type header is not application/json",
			Status:  http.StatusUnsupportedMediaType,
		}
	}

	// Limit request body size to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dst); err != nil {
		return &HTTPError{
			Message: "Request body contains badly-formed JSON",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}

// RespondWithError sends a JSON error response
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(response); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// EnableCORS sets CORS headers for the response
func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// HandleOptions handles preflight OPTIONS requests
func HandleOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		EnableCORS(w)
		w.WriteHeader(http.StatusOK)
		return
	}
}

// package utils

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// )

// // DecodeJSONBody decodes the JSON body of a request into the given interface
// func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
// 	if r.Header.Get("Content-Type") != "application/json" {
// 		return &HTTPError{Message: "Content-Type header is not application/json", Status: http.StatusUnsupportedMediaType}
// 	}

// 	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB max

// 	dec := json.NewDecoder(r.Body)
// 	dec.DisallowUnknownFields()

// 	if err := dec.Decode(&dst); err != nil {
// 		return &HTTPError{Message: "Request body contains badly-formed JSON", Status: http.StatusBadRequest}
// 	}

// 	return nil
// }

// // HTTPError represents an HTTP error
// type HTTPError struct {
// 	Message string
// 	Status  int
// }

// func (e *HTTPError) Error() string {
// 	return e.Message
// }

// // RespondWithError sends a JSON error response
// func RespondWithError(w http.ResponseWriter, code int, message string) {
// 	RespondWithJSON(w, code, map[string]string{"error": message})
// }

// // RespondWithJSON sends a JSON response
// func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
// 	response, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Printf("Error marshaling JSON response: %v", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(response)
// }
