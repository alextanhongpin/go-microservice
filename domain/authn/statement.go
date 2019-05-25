package authn

import "github.com/alextanhongpin/go-microservice/pkg/gostmt"

const (
	_ gostmt.ID = iota
	createToken
	createUser
	deleteToken
	deleteExpiredTokens
	resetPassword
	tokenWithValue
	updateUserPassword
	userWithEmail
)

var statements = gostmt.Raw{
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
	updateUserPassword: `
		UPDATE 	user
		SET 	hashed_password = ?
		WHERE 	id = UUID_TO_BIN(?)
	`,
	createToken: `
		INSERT INTO token (id, token)
		VALUES 	(UUID_TO_BIN(?), HEX(?))
		ON DUPLICATE KEY UPDATE 
			token = token
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
		SELECT 	id, token, created_at
		FROM 	token
		WHERE 	token = UNHEX(?)
		LIMIT 	1
	`,
}
