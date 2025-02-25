package rotas

import (
	"api/src/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type Rota struct {
	URI                string
	Metodo             string
	Funcao             func(http.ResponseWriter, *http.Request)
	RequerAutenticacao bool
}

func Configurar(r *mux.Router) *mux.Router {
	rotas := rotasUsers
	rotas = append(rotas, rotasLogin)
	rotas = append(rotas, rotasPublications...)

	for _, rota := range rotas {
		if rota.RequerAutenticacao {
			r.HandleFunc(rota.URI, middleware.Logger(middleware.Auth(rota.Funcao))).Methods(rota.Metodo)
		} else {
			r.HandleFunc(rota.URI, middleware.Logger(rota.Funcao)).Methods(rota.Metodo)
		}
	}
	return r
}
