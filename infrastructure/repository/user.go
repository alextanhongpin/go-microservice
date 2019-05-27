package repository

import (
	"database/sql"
	"log"

	"github.com/alextanhongpin/go-microservice/domain/user"
	"github.com/alextanhongpin/go-microservice/pkg/gostmt"
	"github.com/alextanhongpin/go-microservice/pkg/mysqluuid"
)

const (
	_ gostmt.ID = iota
	createUser
	changeUserPassword
	userWithEmail
	userWithID
)

var userStmts = gostmt.Raw{
	userWithID: `
		SELECT 	BIN_TO_UUID(id, true), 
			hashed_password
		FROM 	user
		WHERE 	id = UUID_TO_BIN(?, true) 
	`,
	userWithEmail: `
		SELECT 	BIN_TO_UUID(id, true), 
			hashed_password 
		FROM 	user 
		WHERE 	email = ?
	`,
	createUser: `
		INSERT 	INTO user 
			(id, email, hashed_password)
		VALUES 	(UUID_TO_BIN(?, true), ?, ?)
	`,
	changeUserPassword: `
		UPDATE 	user
		SET 	hashed_password = ?
		WHERE 	id = UUID_TO_BIN(?, true)
	`,
}

type User struct {
	stmts gostmt.Statements
}

func NewUser(db *sql.DB) *User {
	stmts, err := gostmt.Prepare(db, userStmts)
	if err != nil {
		log.Fatal("error initializing user statement:", err)
	}
	return &User{stmts}
}

// UserWithEmail returns a User given a valid email.
func (u *User) WithEmail(email string) (user.Entity, error) {
	var usr user.Entity
	err := u.stmts[userWithEmail].QueryRow(email).Scan(
		&usr.ID,
		&usr.HashedPassword,
	)
	return usr, err
}

// UserWithID returns a User given a valid id.
func (u *User) WithID(userID string) (user.Entity, error) {
	var usr user.Entity
	err := u.stmts[userWithID].QueryRow(userID).Scan(
		&usr.ID,
		&usr.HashedPassword,
	)
	return usr, err
}

// Create creates a new User with the given username and password.
func (u *User) Create(username, password string) (user.Entity, error) {
	var usr user.Entity
	// MySQL is using uuid v1.
	usr.ID = mysqluuid.New()
	_, err := u.stmts[createUser].Exec(
		usr.ID,
		username,
		password,
	)
	return usr, err
}

func (u *User) ChangePassword(userID, password string) (bool, error) {
	res, err := u.stmts[changeUserPassword].Exec(password, userID)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}
