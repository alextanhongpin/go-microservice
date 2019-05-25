package authn

import (
	"database/sql"
	"log"
	"time"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/gostmt"
)

type (
	Repository interface {
		// Setter.
		CreateToken(userID, token string) (bool, error)
		CreateUser(username, password string) (User, error)

		DeleteToken(token string) (bool, error)
		DeleteTokens(ttl time.Duration) (int64, error)
		UpdateUserPassword(userID, password string) (bool, error)

		// Getter.
		TokenWithValue(token string) (Token, error)
		UserWithEmail(email string) (User, error)
	}
	// RepositoryImpl implements the RepositoryImpl interface.
	RepositoryImpl struct {
		stmts gostmt.Statements
	}
)

// NewRepositoryImpl returns a new RepositoryImpl.
func NewRepository(db *sql.DB) *RepositoryImpl {
	stmts, err := gostmt.Prepare(db, statements)
	if err != nil {
		log.Fatal(err)
	}
	return &RepositoryImpl{stmts}
}

// UserWithEmail returns a User given a valid email.
func (r *RepositoryImpl) UserWithEmail(email string) (User, error) {
	var user User
	err := r.stmts[userWithEmail].QueryRow(email).Scan(
		&user.ID,
		&user.HashedPassword,
	)
	return user, err
}

// Create creates a new User with the given username and password.
func (r *RepositoryImpl) CreateUser(username, password string) (User, error) {
	var u User
	// MySQL is using uuid v1.
	u.ID = api.NewUUID()
	_, err := r.stmts[createUser].Exec(
		u.ID,
		username,
		password,
	)
	if err != nil {
		return u, err
	}
	return u, err
}

func (r *RepositoryImpl) UpdateUserPassword(userID, password string) (bool, error) {
	res, err := r.stmts[updateUserPassword].Exec(password, userID)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

func (r *RepositoryImpl) CreateToken(userID, token string) (bool, error) {
	res, err := r.stmts[createToken].Exec(userID, token)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

func (r *RepositoryImpl) TokenWithValue(token string) (Token, error) {
	var result Token
	err := r.stmts[tokenWithValue].QueryRow(token).Scan(
		&result.UserID,
		&result.Token,
		&result.CreatedAt,
	)
	return result, err
}

func (r *RepositoryImpl) DeleteToken(token string) (bool, error) {
	res, err := r.stmts[deleteToken].Exec(token)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

func (r *RepositoryImpl) DeleteTokens(ttl time.Duration) (int64, error) {
	minute := int(ttl.Minutes())
	res, err := r.stmts[deleteTokens].Exec(minute)
	if err != nil {
		return -1, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}
