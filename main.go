package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// HomeHandler is the main path
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintln(w, "<h1>Welcome to the Shout Server</h1>")

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
