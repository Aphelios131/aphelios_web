package models

import (
	"database/sql"
	"fmt"
	"strings"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Id            int
	Email         string
	Password_Hash string
}

type UserModel struct {
	DB *sql.DB
}

//user 插入数据
func (um *UserModel) Insert(email, password string) error {
	hash_byte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	email = strings.ToLower(email)
	Password_hash := string(hash_byte)

	stmt := `insert into users (email, password_hash) values ($1,$2)`

	_, err = um.DB.Exec(stmt, email, Password_hash)
	if err != nil {
		return err
	}

	return nil
}

func (um *UserModel) Get(email, password string) (int, error) {
	var id int
	var password_hash []byte
	email = strings.ToLower(email)
	stmt := `select id, password_hash from users where email = $1`
	rows := um.DB.QueryRow(stmt, email)
	err := rows.Scan(&id, &password_hash)
	if err != nil {
		return 0, fmt.Errorf("scanf data err:%w", err)
	}

	err = bcrypt.CompareHashAndPassword(password_hash, []byte(password))
	if err != nil {
		return 0, fmt.Errorf("CompareHashAndPassword err:%w", err)
	}
	return id, nil
}

func (um *UserModel) GetId(id int) (bool, error) {
	var exists bool
	stmt := `select exists(select true from users where id = $1)`
	err := um.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err

}
