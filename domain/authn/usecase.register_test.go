package authn_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/alextanhongpin/go-microservice/domain/authn"
	testdata "github.com/alextanhongpin/go-microservice/testdata"
)

func TestRegisterAccount(t *testing.T) {
	Convey("UseCase: Register account", t, func() {
		db := testdata.NewDB()
		repo := authn.NewRepository(db)
		usecase := authn.NewRegisterUseCase(repo)

		var (
			username = "john.doe@mail.com"
			password = "123456789"
		)

		Convey("Given that the account does not exist", func() {
			Convey("When the user registers a new account", func() {
				_, err := usecase.Register(context.TODO(), authn.RegisterRequest{
					Username: username,
					Password: password,
				})
				Convey("Then the account should be created", func() {
					So(err, ShouldBeNil)
				})
				Convey("And the user cannot register with the same email", func() {
					_, err := usecase.Register(context.TODO(), authn.RegisterRequest{
						Username: username,
						Password: password,
					})
					So(err, ShouldEqual, authn.ErrUserExists)
				})
				Convey("Then the system should store a hashed password", func() {
					user, err := repo.UserWithEmail(username)
					So(err, ShouldBeNil)
					So(user.HashedPassword, ShouldNotEqual, password)
				})
			})
		})
		Reset(func() {
			db.Close()
		})
	})

}
