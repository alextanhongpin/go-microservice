package model

type User struct {
	HashedPassword string `json:"-"`
}
