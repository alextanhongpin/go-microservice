package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/alextanhongpin/go-microservice/domain/token"
	"github.com/alextanhongpin/go-microservice/pkg/gostmt"
)

const (
	_ gostmt.ID = iota
	createToken
	deleteExpiredTokens
	deleteToken
	tokenWithValue
)

var tokenStmts = gostmt.Raw{
	createToken: `
		INSERT INTO token (id, token)
		VALUES 	(UUID_TO_BIN(?, true), UNHEX(?))
		ON DUPLICATE KEY UPDATE 
			token = UNHEX(?)
	`,
	deleteExpiredTokens: `
		DELETE FROM token
		WHERE 	created_at < DATE_SUB(NOW(), INTERVAL ? MINUTE)
	`,
	deleteToken: `
		DELETE FROM token
		WHERE token = UNHEX(?)
	`,
	tokenWithValue: `
		SELECT 	BIN_TO_UUID(id, true), HEX(token), created_at
		FROM 	token
		WHERE 	token = UNHEX(?)
		LIMIT 	1
	`,
}

type Token struct {
	stmts gostmt.Statements
}

func NewToken(db *sql.DB) *Token {
	stmts, err := gostmt.Prepare(db, tokenStmts)
	if err != nil {
		log.Fatal("error initializing token statement:", err)
	}
	return &Token{stmts}
}

func (t *Token) WithValue(tokenString string) (token.Entity, error) {
	var result token.Entity
	err := t.stmts[tokenWithValue].QueryRow(tokenString).Scan(
		&result.UserID,
		&result.Token,
		&result.CreatedAt,
	)
	return result, err
}

func (t *Token) Delete(tkn string) (bool, error) {
	res, err := t.stmts[deleteToken].Exec(tkn)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

func (t *Token) DeleteExpired(ttl time.Duration) (int64, error) {
	minute := int(ttl.Minutes())
	res, err := t.stmts[deleteExpiredTokens].Exec(minute)
	if err != nil {
		return -1, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

func (t *Token) Create(userID, token string) (bool, error) {
	res, err := t.stmts[createToken].Exec(userID, token, token)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}
