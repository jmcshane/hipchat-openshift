package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmcshane/hipchat-openshift/service"
)

// OcCredsHandler Handle post requests that set the current token
type OcCredsHandler struct {
	TokenService *service.TokenService
}

func (oc OcCredsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&oc.TokenService)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		fmt.Fprintf(w, "Test")
	default:
		http.Error(w, fmt.Sprintf("Method %s not supported", r.Method), 404)
	}
}
