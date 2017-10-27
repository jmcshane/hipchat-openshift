package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Server struct {
	*negroni.Negroni
}

func NewServer() *Server {
	s := Server{negroni.Classic()}
	r := http.NewServeMux()
	router := mux.NewRouter()

}
