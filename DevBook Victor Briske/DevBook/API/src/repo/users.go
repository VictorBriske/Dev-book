package repo

import (
	"api/src/models"
	"database/sql"
	"errors"
	"fmt"
)

type Users struct {
	db *sql.DB
}

func NewRepoUsers(db *sql.DB) *Users {
	return &Users{db}
}

func (repo Users) CreateUser(user models.User) (uint64, error) {
	statement, err := repo.db.Prepare("INSERT INTO users (name, nick, email, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(user.Name, user.Nick, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	latestIDInserted, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(latestIDInserted), nil

}

func (repo Users) Search(nameOrNick string) ([]models.User, error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick)

	tuples, err := repo.db.Query(
		"select id, name, nick, email, createdat from users where name LIKE ? or nick like ?",
		nameOrNick, nameOrNick,
	)

	if err != nil {
		return nil, err
	}

	defer tuples.Close()

	var users []models.User

	for tuples.Next() {
		var user models.User

		if err = tuples.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Createdat,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (repo Users) SearchByID(UserID uint64) (models.User, error) {
	tuples, err := repo.db.Query(
		"select id, name, nick,email, createdat from users where id = ?",
		UserID,
	)
	if err != nil {
		return models.User{}, err
	}
	var user models.User
	if tuples.Next() {
		if err = tuples.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Createdat,
		); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (repo Users) UpdateUser(ID uint64, user models.User) error {
	statement, erro := repo.db.Prepare("update users set name = ?, nick = ?, email = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(user.Name, user.Nick, user.Email, ID); erro != nil {
		return erro
	}

	return nil
}

func (repo Users) DeleteUser(ID uint64) error {
	statement, erro := repo.db.Prepare("delete from users where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(ID); erro != nil {
		return erro
	}

	return nil
}

func (repo Users) SearchByEmail(email string) (models.User, error) {
	tuple, err := repo.db.Query("select id, password from users where email = ?", email)
	if err != nil {
		return models.User{}, err
	}
	defer tuple.Close()

	var user models.User

	if tuple.Next() {
		if err = tuple.Scan(&user.ID, &user.Password); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (repo Users) UserExists(userID uint64) (bool, error) {
	var exists bool
	err := repo.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (repo Users) Follow(userID, followerID uint64) error {
	statement, erro := repo.db.Prepare("insert ignore into followers(user_id, follower_id) values (?, ?)")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(userID, followerID); erro != nil {
		return erro
	}
	return nil
}

func (repo Users) StopFollow(userID, followerID uint64) error {
	statement, erro := repo.db.Prepare("delete from followers where user_id = ? and follower_id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(userID, followerID); erro != nil {
		return erro
	}
	return nil
}

func (repo Users) FindAllFollowers(userID uint64) ([]models.User, error) {
	tuples, err := repo.db.Query(
		`SELECT u.id, u.name, u.nick, u.email, u.createdat
		 FROM followers f
		 INNER JOIN users u ON f.follower_id = u.id
		 WHERE f.user_id = ?`, userID,
	)

	if err != nil {
		return nil, err
	}
	defer tuples.Close()

	var users []models.User

	for tuples.Next() {
		var user models.User

		if err = tuples.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Createdat,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, errors.New("não há seguidores")
	}

	return users, nil
}

func (repo Users) FindFollows(follower_id uint64) ([]models.User, error) {
	tuples, err := repo.db.Query(
		`SELECT u.id, u.name, u.nick, u.email, u.createdat
		FROM followers f
		INNER JOIN users u ON f.user_id = u.id
		WHERE f.follower_id = ?`, follower_id,
	)

	if err != nil {
		return nil, err
	}
	defer tuples.Close()

	var users []models.User

	for tuples.Next() {
		var user models.User

		if err = tuples.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Createdat,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, errors.New("este usuário não segue ninguém")
	}

	return users, nil
}

func (repo Users) SearchPassword(userID uint64) (string, error) {
	tuple, err := repo.db.Query("select password from users where id = ?", userID)
	if err != nil {
		return "", err
	}
	defer tuple.Close()

	var user models.User

	if tuple.Next() {
		if err = tuple.Scan(&user.Password); err != nil {
			return "", err
		}
	}

	return user.Password, nil
}

func (repo Users) UpdatePassword(userID uint64, password string) error {
	statement, err := repo.db.Prepare("update users set password = ? where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(password, userID); err != nil {
		return err
	}

	return nil
}
