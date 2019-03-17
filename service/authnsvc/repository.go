package authnsvc

import (
	"database/sql"

	"github.com/alextanhongpin/go-microservice/api"
)

type (
	// Repository represents the data access layer to the User repository.
	Repository interface {
		// Reader.
		WithEmail(email string) (User, error)

		// Writer.
		Create(username, password string) (User, error)
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
func (r *RepositoryImpl) WithEmail(email string) (User, error) {
	stmt := `
		SELECT 
			BIN_TO_UUID(id, true), 
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
func (r *RepositoryImpl) Create(username, password string) (User, error) {
	var u User
	// MySQL is using uuid v1.
	u.ID = api.NewUUID()
	stmt := `
		INSERT INTO user 
			(id, email, hashed_password)
		VALUES (UUID_TO_BIN(?, true), ?, ?)
	`
	_, err := r.db.Exec(stmt,
		u.ID,
		username,
		password,
	)
	if err != nil {
		return u, err
	}
	return u, err
}
