package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasPublications = []Rota{
	{
		URI:                "/publications",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CreatePub,
		RequerAutenticacao: true,
	},
	{
		URI:                "/publications",
		Metodo:             http.MethodGet,
		Funcao:             controllers.SearchPubs,
		RequerAutenticacao: true,
	},
	{
		URI:                "/publications/{publicationId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.SearchPub,
		RequerAutenticacao: true,
	},
	{
		URI:                "/publications/{publicationId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.UpdatePub,
		RequerAutenticacao: true,
	},
	{
		URI:                "/publications/{publicationId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletePub,
		RequerAutenticacao: true,
	},
	{
		URI:                "/users/{userId}/publications",
		Metodo:             http.MethodGet,
		Funcao:             controllers.UserPubs,
		RequerAutenticacao: true,
	},
	{

		URI:                "/publications/{publicationId}/like",
		Metodo:             http.MethodPost,
		Funcao:             controllers.LikePub,
		RequerAutenticacao: true,
	},
	{

		URI:                "/publications/{publicationId}/deslike",
		Metodo:             http.MethodPost,
		Funcao:             controllers.DeslikePub,
		RequerAutenticacao: true,
	},
}
