package repo

import (
	"api/src/models"
	"database/sql"
)

type Users struct {
	db *sql.DB
}

func NewRepoUsers(db *sql.DB) *Users {
	return &Users{db}
}

func (repo Users) CreateUser(user models.User) (uint64, error) {
	statement, erro := repo.db.Prepare("insert into users (name, nick, email, password) values (?, ?, ?, ?)")
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	result, erro := statement.Exec(user.Name, user.Nick, user.Email, user.Password)
	if erro != nil {
		return 0, erro
	}

	LatestIDInserted, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(LatestIDInserted), nil
}
