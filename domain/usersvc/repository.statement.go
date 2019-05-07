package usersvc

import "github.com/alextanhongpin/go-microservice/pkg/gostmt"

const (
	_ gostmt.ID = iota
	withID
)

var statements = gostmt.Raw{
	withID: `
		SELECT 	BIN_TO_UUID(id, true) AS uuid,
			name,
			picture,
			created_at,
			birthdate
		FROM 	user
		WHERE 	id = UUID_TO_BIN(?, true)
		LIMIT   1
	`,
}
