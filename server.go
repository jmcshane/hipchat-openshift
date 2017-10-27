package main

import (
	"log"
	"net/http"
)

//TokenService share a token between the creds handler and hipchat handler
type TokenService struct {
	Token string
}

func main() {
	tokenService := TokenService{}
	hipchatHandler := HipchatHandler{tokenService: &tokenService}
	ocHandler := OcCredsHandler{tokenService: &tokenService}
	http.Handle("/creds", ocHandler)
	http.Handle("/", hipchatHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
