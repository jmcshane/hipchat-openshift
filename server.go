package main

import (
	"log"
	"net/http"

	"github.com/jmcshane/hipchat-openshift/controllers"
	"github.com/jmcshane/hipchat-openshift/service"
)

func main() {
	tokenService := service.NewTokenService()
	tokenService.Token = "t0O84cpOfa75gHUoksbehRXoh4_z_ZbfPFGry-r1XeI"
	hipchatHandler := controllers.NewHipchatHandler(tokenService)
	ocHandler := controllers.NewOcCredsHandler(tokenService)
	http.Handle("/creds", ocHandler)
	http.Handle("/", hipchatHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
