package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// APIResponse is the response to the request
type APIResponse struct {
	Code    int
	Type    string
	Message interface{}
}

// HomeHandler is the main path
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := APIResponse{
		Code:    http.StatusOK,
		Type:    "unknown",
		Message: "Welcome to the Bid Bang",
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/user", UserRegistrationHandler)
	http.Handle("/", r)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
