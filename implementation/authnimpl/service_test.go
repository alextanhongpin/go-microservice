package authnimpl_test

import (
	"testing"

	"github.com/alextanhongpin/go-microservice/domain/user"
	"github.com/alextanhongpin/go-microservice/implementation/authnimpl"
)

func TestAuthnService(t *testing.T) {
	svc := authnimpl.NewService()

	t.Run("should return error when hashing empty password", func(t *testing.T) {
		t.Parallel()
		_, err := svc.HashPassword("")
		if err != nil && err != authnimpl.ErrPasswordRequired {
			t.Fatalf("want %v, got %v", authnimpl.ErrPasswordRequired, err)
		}
	})

	hashed, err := svc.HashPassword("hello world")
	if err != nil {
		t.Fatal(err)
	}
	u := user.Entity{
		HashedPassword: hashed,
	}
	t.Run("should return true when password match", func(t *testing.T) {
		t.Parallel()
		ok := svc.ComparePassword(u, "hello world")
		if !ok {
			t.Fatalf("want %t, got %t", true, ok)
		}
	})
	t.Run("should return false when comparing against an empty string", func(t *testing.T) {
		t.Parallel()
		ok := svc.ComparePassword(u, "")
		if ok {
			t.Fatalf("want %t, got %t", false, ok)
		}
	})
	t.Run("should return false when password do not match", func(t *testing.T) {
		t.Parallel()
		ok := svc.ComparePassword(u, "hello")
		if ok {
			t.Fatalf("want %t, got %t", false, ok)
		}
	})
	t.Run("should match case-sensitive password", func(t *testing.T) {
		t.Parallel()
		ok := svc.ComparePassword(u, "Hello World")
		if ok {
			t.Fatalf("want %t, got %t", false, ok)
		}
	})
}
