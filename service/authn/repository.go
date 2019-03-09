package authn

import (
	"database/sql"

	uuid "github.com/satori/go.uuid"
)

type (
	// Repository represents the data access layer to the User repository.
	Repository interface {
		GetUser(email string) (User, error)
		CreateUser(username, password string) (User, error)
	}
	// RepositoryImpl implements the Repository interface.
	RepositoryImpl struct {
		db *sql.DB
	}
)

// NewRepository returns a new Repository.
func NewRepository(db *sql.DB) *RepositoryImpl {
	return &RepositoryImpl{db}
}

// GetUser returns a User given a valid email.
func (r *RepositoryImpl) GetUser(email string) (User, error) {
	stmt := `
		SELECT 
			id, 
			hashed_password 
		FROM user 
		WHERE email = ?
	`
	var user User
	err := r.db.QueryRow(stmt, email).Scan(
		&user.ID,
		&user.HashedPassword,
	)
	return user, err
}

// CreateUser creates a new User with the given username and password.
func (r *RepositoryImpl) CreateUser(username, password string) (User, error) {
	var user User
	// MySQL is using uuid v1.
	user.ID = uuid.Must(uuid.NewV1()).String()
	stmt := `
		INSERT INTO user 
			(id, email, hashed_password)
		VALUES (UUID_TO_BIN(?, true), ?, ?)
	`
	_, err := r.db.Exec(stmt,
		user.ID,
		username,
		password,
	)
	if err != nil {
		return user, err
	}
	return user, err
}
