package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Option struct {
	User string `envconfig:"DB_USER" required:"true"`
	Pass string `envconfig:"DB_PASS" required:"true"`
	Host string `envconfig:"DB_HOST" required:"true"`
	Name string `envconfig:"DB_NAME" required:"true"`
}

func NewOption() (Option, error) {
	var opt Option
	err := envconfig.Process("", &opt)
	return opt, err
}

func New(opt Option) (*sql.DB, error) {
	cfg := mysql.Config{
		User:      opt.User,
		Passwd:    opt.Pass,
		Addr:      opt.Host,
		DBName:    opt.Name,
		ParseTime: true,
		Params:    map[string]string{"charset": "utf8"},
		Collation: "utf8mb4_unicode_ci",
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewProduction() *sql.DB {
	opt, err := NewOption()
	if err != nil {
		log.Fatal(err.Error(), zap.Error(err))
	}
	db, err := New(opt)
	if err != nil {
		db.Close()
		log.Fatal(err.Error(), zap.Error(err))
	}
	for i := 0; i < 3; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		db.Close()
		log.Fatal(err.Error(), zap.Error(err))
	}
	return db
}
