package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type requestHandler struct{}

func (handler requestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		cmd := exec.Command("oc", "get", "pods")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s", out.String())
	default:
		http.Error(w, fmt.Sprintf("Method %s not supported", r.Method), 404)
	}
}

func main() {
	handler := requestHandler{}
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
