package database_test

import (
	"os"
	"testing"

	"github.com/alextanhongpin/go-microservice/database"
)

func TestReadConfig(t *testing.T) {
	var (
		user = "user"
		name = "test-db"
		pass = "123456"
		host = "host"
	)

	os.Setenv("DB_NAME", name)
	os.Setenv("DB_USER", user)
	os.Setenv("DB_PASS", pass)
	os.Setenv("DB_HOST", host)

	opt, err := database.NewOption()
	if err != nil {
		t.Fatal(err)
	}
	if v := opt.Name; v != name {
		t.Fatalf("expected %v, got %v", name, v)
	}
	if v := opt.User; v != user {
		t.Fatalf("expected %v, got %v", user, v)
	}
}
