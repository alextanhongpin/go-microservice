package user

import (
	"database/sql"
	"log"

	"github.com/alextanhongpin/go-microservice/pkg/gostmt"
)

type (
	Repository struct {
		stmts gostmt.Statements
	}
)

func NewRepository(db *sql.DB) *Repository {
	stmts, err := gostmt.Prepare(db, statements)
	if err != nil {
		log.Fatal(err)
	}
	return &Repository{stmts}
}

func (r *Repository) WithID(id string) (User, error) {
	var u User
	err := r.stmts[withID].QueryRow(id).Scan(
		&u.ID,
		&u.Name,
		&u.Picture,
		&u.CreatedAt,
		&u.BirthDate,
	)
	return u, err
}

func (r *Repository) BelongingToPage() ([]User, error) {
	return nil, nil
}
