package models

import (
	"api/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	Createdat time.Time `json:"createdat,omitempty"`
}

func (user *User) Prepare(step string) error {
	if erro := user.validate(step); erro != nil {
		return erro
	}

	if err := user.format(step); err != nil {
		return err
	}
	return nil
}

func (user *User) validate(step string) error {
	if user.Name == "" {
		return errors.New("o nome é obrigatório")
	}

	if user.Nick == "" {
		return errors.New("o nick é obrigatório")
	}

	if user.Email == "" {
		return errors.New("o email é obrigatório")
	}

	if erro := checkmail.ValidateFormat(user.Email); erro != nil {
		return errors.New("o email inserido é inválido")
	}

	if step == "cadastro" && user.Password == "" {
		return errors.New("a senha é obrigatória")
	}

	return nil
}

func (user *User) format(step string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Nick = strings.TrimSpace(user.Nick)
	user.Email = strings.TrimSpace(user.Email)

	if step == "cadastro" {
		hashedpassword, err := security.Hash(user.Password)
		if err != nil {
			return err
		}

		user.Password = string(hashedpassword)

	}

	return nil
}
