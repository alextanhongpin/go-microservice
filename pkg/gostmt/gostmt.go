package gostmt

import "database/sql"

// ID is the id of the statement.
type ID uint

// Raw represents the map of unprepared statements.
type Raw map[ID]string

// Statements represent the prepared statement.
type Statements map[ID]*sql.Stmt

// Prepare takes in the db connection and prepares the raw statement.
func Prepare(db *sql.DB, stmts Raw) (Statements, error) {
	res := make(Statements)
	for id, stmt := range stmts {
		var err error
		res[id], err = db.Prepare(stmt)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
