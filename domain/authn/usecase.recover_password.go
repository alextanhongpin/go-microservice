package authn

import (
	"context"
	"log"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/alextanhongpin/go-microservice/pkg/govalidator"
)

type (
	RecoverPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}
	RecoverPasswordResponse struct {
		Success bool   `json:"success"`
		Token   string `json:"-"`
	}
)

type (
	recoverPasswordUseCase interface {
		RecoverPassword(ctx context.Context, req RecoverPasswordRequest) (*RecoverPasswordResponse, error)
	}
	recoverPasswordRepository interface {
		UserWithEmail(email string) (User, error)
		CreateToken(userID, token string) (bool, error)
		DeleteExpiredTokens(ttl time.Duration) (int64, error)
	}
)

type RecoverPasswordUseCase struct {
	repo     recoverPasswordRepository
	tokenTTL time.Duration
}

func NewRecoverPasswordUseCase(
	repo recoverPasswordRepository,
	tokenTTL time.Duration,
) (*RecoverPasswordUseCase, func()) {
	usecase := &RecoverPasswordUseCase{
		repo:     repo,
		tokenTTL: tokenTTL,
	}
	cancel := usecase.init()
	return usecase, cancel
}

func (r *RecoverPasswordUseCase) RecoverPassword(ctx context.Context, req RecoverPasswordRequest) (*RecoverPasswordResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, err
	}
	user, err := r.repo.UserWithEmail(req.Email)
	if err != nil {
		// Obfuscate error to prevent attacker from guessing if the
		// user exists.
		return nil, ErrInvalidRequest
	}

	// We only store the sha256 hashed token in the database. The database
	// should only contain the hash, and the created_at time. The
	// information in the database should not be tied back to the users in
	// any way (email, user id foreign key etc).
	// ? Hash the token with the email?
	var (
		userID      = user.ID
		token       = uuid.Must(uuid.NewV4()).String()
		tokenHashed = hashToken(token)
	)
	success, err := r.repo.CreateToken(userID, tokenHashed)
	return &RecoverPasswordResponse{
		Success: success,
		// Return the unhashed token to the user in the email link.
		Token: token,
	}, err
}

// periodic job to clear the tokens in the database.
func (r *RecoverPasswordUseCase) init() func() {
	var once sync.Once
	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan struct{})
	go func() {
		defer wg.Done()

		t := time.NewTicker(2 * r.tokenTTL)
		defer t.Stop()

		for {
			select {
			case <-done:
				return
			case <-t.C:
				if count, err := r.repo.DeleteExpiredTokens(r.tokenTTL); err != nil {
					log.Println("error deleting tokens:", err)
				} else {
					log.Println("tokens deleted:", count)
				}
			}
		}
	}()
	return func() {
		once.Do(func() {
			close(done)
			wg.Wait()
		})
	}
}
