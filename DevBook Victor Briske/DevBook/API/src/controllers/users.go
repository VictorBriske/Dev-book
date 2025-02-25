package controllers

import (
	db "api/src/DataBase"
	"api/src/auth"
	"api/src/models"
	"api/src/repo"
	"api/src/responses"
	"api/src/security"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	bodyRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	defer r.Body.Close()

	var user models.User

	erro = json.Unmarshal(bodyRequest, &user)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	erro = user.Prepare("cadastro")
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	dbConn, erro := db.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer dbConn.Close()

	repo := repo.NewRepoUsers(dbConn)

	user.ID, erro = repo.CreateUser(user)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	nameOrNick := strings.ToLower(r.URL.Query().Get("user"))

	db, erro := db.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repo := repo.NewRepoUsers(db)
	users, err := repo.Search(nameOrNick)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, erro := db.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repo := repo.NewRepoUsers(db)

	exists, err := repo.UserExists(userID)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if !exists {
		responses.Erro(w, http.StatusNotFound, errors.New("usuário não encontrado"))
		return
	}

	user, err := repo.SearchByID(userID)

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	TokenUserID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.Erro(w, http.StatusUnauthorized, err)
		return
	}

	if userID != TokenUserID {
		responses.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar um usuário que não seja o seu"))
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(requestBody, &user); err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	if err = user.Prepare("edicao"); err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, erro := db.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repo := repo.NewRepoUsers(db)

	exists, err := repo.UserExists(userID)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if !exists {
		responses.Erro(w, http.StatusNotFound, errors.New("usuário não encontrado"))
		return
	}
	err = repo.UpdateUser(userID, user)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	TokenUserID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.Erro(w, http.StatusUnauthorized, err)
		return
	}

	if userID != TokenUserID {
		responses.Erro(w, http.StatusForbidden, errors.New("não é possível deletar um usuário que não seja o seu"))
		return
	}

	repo := repo.NewRepoUsers(db)

	exists, err := repo.UserExists(userID)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if !exists {
		responses.Erro(w, http.StatusNotFound, errors.New("usuário não encontrado"))
		return
	}
	err = repo.DeleteUser(userID)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, erro := auth.ExtractUserID(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	userID, erro := strconv.ParseUint(params["userID"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if followerID == userID {
		responses.Erro(w, http.StatusForbidden, errors.New("não é possível seguir você mesmo"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repo.NewRepoUsers(db)
	if erro = repo.Follow(userID, followerID); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, erro)
}

func StopFollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, erro := auth.ExtractUserID(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	userID, erro := strconv.ParseUint(params["userID"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if followerID == userID {
		responses.Erro(w, http.StatusForbidden, errors.New("não é possível parar de seguir você mesmo"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repo.NewRepoUsers(db)
	if erro = repo.StopFollow(userID, followerID); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func FindAllFollowers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	var user []models.User
	repo := repo.NewRepoUsers(db)

	user, err = repo.FindAllFollowers(userID)

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)

}

func FindAllFollows(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	follower_id, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	var user []models.User
	repo := repo.NewRepoUsers(db)

	user, err = repo.FindFollows(follower_id)

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)

}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	TokenUserID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	if TokenUserID != userID {
		responses.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar a senha de um usuário que não seja o seu"))
		return
	}

	bodyRequest, _ := io.ReadAll(r.Body)

	var password models.Password
	if err = json.Unmarshal(bodyRequest, &password); err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repo.NewRepoUsers(db)

	BankPassword, err := repo.SearchPassword(userID)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.Verify(BankPassword, password.Actual); err != nil {
		responses.Erro(w, http.StatusUnauthorized, errors.New("senha incorreta"))
		return
	}

	hashedpassword, err := security.Hash(password.New)

	if err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	if err = repo.UpdatePassword(userID, string(hashedpassword)); err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
