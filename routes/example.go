package routes

import (
	"fmt"
	"net/http"
)

func loadExampleMux() *http.ServeMux {
	exampleMux := http.NewServeMux()

	exampleMux.HandleFunc("POST /protected", protected)

	return exampleMux
}

func protected(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	fmt.Fprintf(w, "Welcome, %s!\n", username)
}
