package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasUsers = []Rota{
	{
		URI:                "/users",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CreateUser,
		RequerAutenticacao: false,
	},

	{
		URI:                "/users/{userID}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.GetUserByID,
		RequerAutenticacao: false,
	},

	{
		URI:                "/users",
		Metodo:             http.MethodGet,
		Funcao:             controllers.GetUser,
		RequerAutenticacao: true,
	},

	{
		URI:                "/users/{userID}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.UpdateUser,
		RequerAutenticacao: true,
	},

	{
		URI:                "/users/{userID}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeleteUser,
		RequerAutenticacao: true,
	},
	{
		URI:                "/users/{userID}/follow",
		Metodo:             http.MethodPost,
		Funcao:             controllers.FollowUser,
		RequerAutenticacao: true,
	},
	{
		URI:                "/users/{userID}/stop-follow",
		Metodo:             http.MethodPost,
		Funcao:             controllers.StopFollowUser,
		RequerAutenticacao: true,
	},
	{
		URI:                "/users/{userID}/find-all-followers",
		Metodo:             http.MethodGet,
		Funcao:             controllers.FindAllFollowers,
		RequerAutenticacao: true,
	},
	{
		URI:                "/users/{userID}/findfollows",
		Metodo:             http.MethodGet,
		Funcao:             controllers.FindAllFollows,
		RequerAutenticacao: true,
	},

	{
		URI:                "/users/{userID}/forgotpassword",
		Metodo:             http.MethodPost,
		Funcao:             controllers.ForgotPassword,
		RequerAutenticacao: true,
	},
}
