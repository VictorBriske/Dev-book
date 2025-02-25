package controllers

import (
	db "api/src/DataBase"
	"api/src/auth"
	"api/src/models"
	"api/src/repo"
	"api/src/responses"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreatePub(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.Erro(w, http.StatusUnauthorized, err)
	}

	bodyRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	defer r.Body.Close()

	var pub models.Publication

	erro = json.Unmarshal(bodyRequest, &pub)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	pub.AuthorID = userID

	if erro = pub.Prepare(); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}
	dbConn, erro := db.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer dbConn.Close()

	repo := repo.NewRepoPubs(dbConn)

	pub.ID, erro = repo.CreatePub(pub)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, pub)
}

func SearchPubs(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.Erro(w, http.StatusUnauthorized, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repo.NewRepoPubs(db)

	pubs, err := repo.Search(userID)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, pubs)
}

func SearchPub(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	pubID, err := strconv.ParseUint(params["publicationId"], 10, 64)
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

	repo := repo.NewRepoPubs(db)

	pub, err := repo.SearchByID(pubID)

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, pub)

}

func UpdatePub(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.Erro(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)

	pubID, err := strconv.ParseUint(params["publicationId"], 10, 64)
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

	repo := repo.NewRepoPubs(db)
	bankPubSave, err := repo.SearchByID(pubID)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if bankPubSave.AuthorID != userID {
		responses.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar uma publicação que não seja sua"))
		return
	}

	bodyRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	defer r.Body.Close()

	var publication models.Publication

	if err = json.Unmarshal(bodyRequest, &publication); err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	if err = publication.Prepare(); err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
	}

	if err = repo.Update(pubID, publication); err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusNoContent, nil)
}

func DeletePub(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.Erro(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)

	pubID, err := strconv.ParseUint(params["publicationId"], 10, 64)
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

	repo := repo.NewRepoPubs(db)
	bankPubSave, err := repo.SearchByID(pubID)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if bankPubSave.AuthorID != userID {
		responses.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar uma publicação que não seja sua"))
		return
	}

	if err = repo.Delete(pubID); err != nil {
		responses.Erro(w, http.StatusNoContent, nil)
	}
}

func UserPubs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userId"], 10, 64)
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

	repo := repo.NewRepoPubs(db)
	publications, err := repo.SearchByUser(userID)

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, publications)
}

func LikePub(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	pubID, err := strconv.ParseUint(params["publicationId"], 10, 64)
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

	repo := repo.NewRepoPubs(db)
	if err := repo.Like(pubID); err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, err)

}

func DeslikePub(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	pubID, err := strconv.ParseUint(params["publicationId"], 10, 64)
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

	repo := repo.NewRepoPubs(db)
	if err := repo.Deslike(pubID); err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, err)
}
