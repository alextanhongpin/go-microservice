package authsvc

import (
	"database/sql"
	"strconv"

	"github.com/alextanhongpin/go-microservice/model"
)

type (
	Repository interface {
		GetUser(username string) (model.User, error)
		CreateUser(username, password string) (model.User, error)
	}
	RepositoryImpl struct {
		db *sql.DB
	}
)

func NewRepository(db *sql.DB) *RepositoryImpl {
	return &RepositoryImpl{db}
}

func (r *RepositoryImpl) GetUser(username string) (model.User, error) {
	stmt := `
		SELECT 
			id, 
			hashed_password 
		FROM user 
		WHERE email = ?
	`
	var user model.User
	err := r.db.QueryRow(stmt, username).Scan(
		&user.ID,
		&user.HashedPassword,
	)
	return user, err
}

func (r *RepositoryImpl) CreateUser(username, password string) (model.User, error) {
	var user model.User
	stmt := `
		INSERT INTO user (email, hashed_password)
		VALUES (?, ?)
	`
	rows, err := r.db.Exec(stmt, username, password)
	if err != nil {
		return user, err
	}
	id, err := rows.LastInsertId()
	if err != nil {
		return user, err
	}
	user.ID = strconv.FormatInt(id, 10)
	return user, err
}
