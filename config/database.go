package config

type Database struct {
	User string `envconfig:"DB_USER" required:"true"`
	Pass string `envconfig:"DB_PASS" required:"true"`
	Host string `envconfig:"DB_HOST" required:"true"`
	Name string `envconfig:"DB_NAME" required:"true"`
}
