package authnsvc

import (
	"database/sql"
	"log"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/gostmt"
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
		stmts gostmt.Statements
	}
)

// NewRepository returns a new Repository.
func NewRepository(db *sql.DB) *RepositoryImpl {
	stmts, err := gostmt.Prepare(db, statements)
	if err != nil {
		log.Fatal(err)
	}
	return &RepositoryImpl{stmts}
}

// GetUser returns a User given a valid email.
func (r *RepositoryImpl) WithEmail(email string) (User, error) {
	var user User
	err := r.stmts[withEmailStmt].QueryRow(email).Scan(
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
	_, err := r.stmts[createStmt].Exec(
		u.ID,
		username,
		password,
	)
	if err != nil {
		return u, err
	}
	return u, err
}
