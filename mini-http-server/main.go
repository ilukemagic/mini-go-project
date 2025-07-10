package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	CreateAt string `json:"create_at,omitempty"`
}

// logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

// custom error handler
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	switch status {
	case http.StatusNotFound:
		fmt.Fprintf(w, "404 Not Found")
	case http.StatusInternalServerError:
		fmt.Fprintf(w, "500 Internal Server Error")
	default:
		fmt.Fprintf(w, "An error occurred")
	}

	log.Printf("Error %s %s: %d", r.Method, r.URL.Path, status)
}

// json response
func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// parse the request body
func parseRequestBody(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.Unmarshal(body, v)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	router := http.NewServeMux()

	// home page
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// Only trigger 404 for non-root paths
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page")
	})

	// about page
	router.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the about page")
	})

	// user api
	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// mock user data
			users := []User{
				{ID: 1, Name: "John"},
				{ID: 2, Name: "Jane"},
				{ID: 3, Name: "Jim"},
				{ID: 4, Name: "Jill"},
			}
			responseWithJSON(w, http.StatusOK, users)
			return
		}

		if r.Method == http.MethodPost {
			var newUser User
			if err := parseRequestBody(r, &newUser); err != nil {
				log.Printf("Error decoding JSON: %v", err)
				errorHandler(w, r, http.StatusBadRequest)
				return
			}

			// set create time
			newUser.CreateAt = time.Now().Format(time.RFC3339)

			// response new created user
			responseWithJSON(w, http.StatusCreated, newUser)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// server
	server := &http.Server{
		Addr:    ":8000",
		Handler: loggingMiddleware(router),
	}

	log.Println("Server is running on port 8000")

	// start server
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
