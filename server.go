package main

import (
	"log"
	"net/http"

	"github.com/jmcshane/hipchat-openshift/controllers"
	"github.com/jmcshane/hipchat-openshift/service"
)

func main() {
	tokenService := service.TokenService{}
	tokenService.Token = "t0O84cpOfa75gHUoksbehRXoh4_z_ZbfPFGry-r1XeI"
	hipchatHandler := controllers.HipchatHandler{TokenService: &tokenService}
	ocHandler := controllers.OcCredsHandler{TokenService: &tokenService}
	http.Handle("/creds", ocHandler)
	http.Handle("/", hipchatHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
