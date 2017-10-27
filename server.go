package main

import (
	"bytes"
    "fmt"
    "log"
    "net/http"
    "os/exec"
)

func handler(w http.ResponseWriter, r *http.Request) {
    cmd := exec.Command("oc", "get", "pods")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", out.String())
}

func main() {
    http.HandleFunc("/foo", handler)    
    log.Fatal(http.ListenAndServe(":8080", nil))    
}