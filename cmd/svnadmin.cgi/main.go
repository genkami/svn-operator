package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	if err := cgi.Serve(RootHandler()); err != nil {
		writeDefaultErrorResponse(err)
	}
}

func RootHandler() http.Handler {
	scriptName := os.Getenv("SCRIPT_NAME")
	r := mux.NewRouter()
	r.HandleFunc(scriptName+"/users", UsersHandler).Methods(http.MethodPost)
	r.HandleFunc(scriptName+"/repo/{name:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?}", RepoHandler).Methods(http.MethodPost)
	return r
}

func RepoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "Path: %s\n", r.URL.Path)
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "Path: %s\n", r.URL.Path)
	fmt.Fprintf(w, "Repo Name: %s\n", vars["name"])
}

func writeDefaultErrorResponse(err error) {
	fmt.Printf("Status: 500\n")
	fmt.Printf("Content-Type: text/plain\n")
	fmt.Printf("\n")
	fmt.Printf("Error: %s", err.Error())
}
