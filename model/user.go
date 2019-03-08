package model

type User struct {
	ID             string `json:"-"`
	HashedPassword string `json:"-"`
}
