package config_test

import (
	"os"
	"testing"

	"github.com/alextanhongpin/go-microservice/config"
	"github.com/alextanhongpin/go-microservice/pkg/str"
	"github.com/stretchr/testify/assert"
)

func TestLoadNestedConfig(t *testing.T) {
	assert := assert.New(t)
	var (
		user = str.Rand(8)
		pass = str.Rand(8)
		host = str.Rand(8)
		name = str.Rand(8)
	)
	os.Setenv("DB_USER", user)
	os.Setenv("DB_PASS", pass)
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_NAME", name)

	cfg := config.New()
	assert.Equal(user, cfg.Database.User)
	assert.Equal(pass, cfg.Database.Pass)
	assert.Equal(host, cfg.Database.Host)
	assert.Equal(name, cfg.Database.Name)
}
