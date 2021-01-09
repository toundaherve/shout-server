package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
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

var dbURL = "postgresql://localhost/bidbang?user=herves&password=1981"

func main() {
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/user", NewUserRegistrationHandler(conn)).Methods("POST")
	http.Handle("/", r)

	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
