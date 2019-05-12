package authn

import "github.com/alextanhongpin/go-microservice/pkg/gostmt"

const (
	_ gostmt.ID = iota
	withEmailStmt
	createStmt
)

var statements = gostmt.Raw{
	withEmailStmt: `
		SELECT 	BIN_TO_UUID(id, true), 
			hashed_password 
		FROM 	user 
		WHERE 	email = ?
	`,
	createStmt: `
		INSERT 	INTO user 
			(id, email, hashed_password)
		VALUES 	(UUID_TO_BIN(?, true), ?, ?)
	`,
}
