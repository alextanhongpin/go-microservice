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
	userWithID
)

var statements = gostmt.Raw{
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
	updateUserPassword: `
		UPDATE 	user
		SET 	hashed_password = ?
		WHERE 	id = UUID_TO_BIN(?, true)
	`,
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
