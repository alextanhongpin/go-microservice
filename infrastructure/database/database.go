package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
)

type Option struct {
	User string
	Pass string
	Host string
	Name string
}

func (o Option) ConnectionString() string {
	cfg := mysql.Config{
		User:      o.User,
		Passwd:    o.Pass,
		Addr:      o.Host,
		DBName:    o.Name,
		ParseTime: true,
		Params:    map[string]string{"charset": "utf8"},
		Collation: "utf8mb4_unicode_ci",
		// Required for mysql:8.0.0 and above.
		AllowNativePasswords: true,
	}
	return cfg.FormatDSN()
}

func New(opt Option) (*sql.DB, error) {
	db, err := sql.Open("mysql", opt.ConnectionString())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewProduction(opt Option) *sql.DB {
	db, err := New(opt)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("dbError: %+v, retry=%d\n", err, i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	return db
}

func IsNotFound(err error) bool {
	return err == sql.ErrNoRows
}

func IsDuplicateEntry(err error) bool {
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		return mysqlError.Number == mysqlerr.ER_DUP_ENTRY
	}
	return false
}
