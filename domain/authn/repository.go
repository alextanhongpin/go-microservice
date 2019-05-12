package authn

import (
	"database/sql"
	"log"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/gostmt"
)

type (
	// Repository implements the Repository interface.
	Repository struct {
		stmts gostmt.Statements
	}
)

// NewRepository returns a new Repository.
func NewRepository(db *sql.DB) *Repository {
	stmts, err := gostmt.Prepare(db, statements)
	if err != nil {
		log.Fatal(err)
	}
	return &Repository{stmts}
}

// WithEmail returns a User given a valid email.
func (r *Repository) WithEmail(email string) (User, error) {
	var user User
	err := r.stmts[withEmailStmt].QueryRow(email).Scan(
		&user.ID,
		&user.HashedPassword,
	)
	return user, err
}

// Create creates a new User with the given username and password.
func (r *Repository) Create(username, password string) (User, error) {
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
