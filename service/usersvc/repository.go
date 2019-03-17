package usersvc

import "database/sql"

type (
	Repository interface {
		WithID(id string) (User, error)
	}
	RepositoryImpl struct {
		db *sql.DB
	}
)

func NewRepository(db *sql.DB) *RepositoryImpl {
	return &RepositoryImpl{db}
}

func (r *RepositoryImpl) WithID(id string) (User, error) {
	var u User
	stmt := `
		SELECT
			BIN_TO_UUID(id, true) AS uuid,
			name,
			picture,
			created_at,
			birthdate
		FROM 	user
		WHERE 	id = UUID_TO_BIN(?, true)
		LIMIT   1
	`
	err := r.db.QueryRow(stmt, id).Scan(
		&u.ID,
		&u.Name,
		&u.Picture,
		&u.CreatedAt,
		&u.BirthDate,
	)
	return u, err
}

func (r *RepositoryImpl) BelongingToPage() ([]User, error) {
	return nil, nil
}
