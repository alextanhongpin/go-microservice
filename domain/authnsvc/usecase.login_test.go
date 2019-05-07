package authnsvc_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/alextanhongpin/go-microservice/domain/authnsvc"
	"github.com/alextanhongpin/go-microservice/pkg/govalidator"

	"github.com/alextanhongpin/passwd"
)

func TestLoginRequest(t *testing.T) {
	tests := []struct {
		email, password string
		isErr           bool
	}{
		{"a.b@mail.com", "12345678", false},
		{"", "", true},
		{"a.b@mail.com", "", true},
		{"a.b@mail.com", "1234567", true},
		{"", "1234567", true},
		{"", "12345678", true},
	}
	for _, tt := range tests {
		err := govalidator.Validate.Struct(authnsvc.LoginRequest{
			Username: tt.email,
			Password: tt.password,
		})
		if isErr := err != nil; isErr != tt.isErr {
			t.Fatal(err)
		}
	}
}

func TestLogin(t *testing.T) {
	Convey("UseCase Login", t, func() {
		var (
			email    = "john.doe@mail.com"
			password = "12345678"
			userID   = "1"
		)
		// Build usecase.
		Convey("given a registered User", func() {
			// A registered user must have a hashed password.
			hashedPwd, err := passwd.Hash(password)
			So(err, ShouldBeNil)
			user := authnsvc.User{ID: userID, HashedPassword: hashedPwd}
			repo := newRepository(user, nil)
			usecase := authnsvc.NewLoginUseCase(repo)

			Convey("when the User enters a valid email and password", func() {
				req := authnsvc.LoginRequest{
					Username: email,
					Password: password,
				}
				Convey("then the User should receive an access token", func() {
					res, err := usecase.Login(context.Background(), req)
					So(err, ShouldBeNil)
					So(res, ShouldNotBeNil)
					So(repo.invoked, ShouldBeTrue)
					So(repo.invokedCount, ShouldEqual, 1)
					So(res.User.ID, ShouldEqual, userID)
				})
			})
			Convey("when the User enters an invalid password (len < 8)", func() {
				req := authnsvc.LoginRequest{
					Username: email,
					Password: "",
				}
				So(len(req.Password), ShouldBeLessThan, 8)
				Convey("then the system should respond with a validation error", func() {
					res, err := usecase.Login(context.Background(), req)
					So(res, ShouldBeNil)
					So(err.Error(), ShouldContainSubstring, "Password")
					So(err.Error(), ShouldContainSubstring, "required")
					So(repo.invoked, ShouldBeFalse)
					So(repo.invokedCount, ShouldEqual, 0)
				})
			})
			Convey("when the User enters the wrong password", func() {
				req := authnsvc.LoginRequest{
					Username: email,
					Password: "87654321",
				}
				So(req.Password, ShouldNotEqual, password)
				Convey("then the system should respond with an error", func() {
					res, err := usecase.Login(context.Background(), req)
					So(res, ShouldBeNil)
					So(err.Error(), ShouldContainSubstring, "password do not match")
					So(repo.invoked, ShouldBeTrue)
					So(repo.invokedCount, ShouldEqual, 1)
				})
			})
			// TODO: When the user enters the wrong password three times.
		})
		Convey("given a unregistered User", func() {
			var errUserDoesNotExist = errors.New("user does not exist")
			repo := newRepository(authnsvc.User{}, errUserDoesNotExist)
			usecase := authnsvc.NewLoginUseCase(repo)
			Convey("when the User enters a fake username and password", func() {
				req := authnsvc.LoginRequest{
					Username: "jane.doe@mail.com",
					Password: "xyzabc123",
				}
				Convey("then the system should respond with an error", func() {
					res, err := usecase.Login(context.Background(), req)
					So(res, ShouldBeNil)
					So(err.Error(), ShouldContainSubstring, errUserDoesNotExist.Error())
					So(repo.invoked, ShouldBeTrue)
					So(repo.invokedCount, ShouldEqual, 1)
				})
			})
		})
	})
}

// Helpers.

type repository struct {
	user         authnsvc.User
	err          error
	invoked      bool
	invokedCount int
}

func newRepository(user authnsvc.User, err error) *repository {
	return &repository{user, err, false, 0}
}

func (r *repository) WithEmail(email string) (authnsvc.User, error) {
	defer func() {
		r.invoked = true
		r.invokedCount++
	}()
	return r.user, r.err
}
