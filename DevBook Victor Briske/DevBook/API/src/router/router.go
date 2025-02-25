package router

import (
	"api/src/router/rotas"

	"github.com/gorilla/mux"
)

func Gerar() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	return rotas.Configurar(r)
}
