package tokenimpl

import (
	"github.com/alextanhongpin/go-microservice/pkg/gostrings"
	"github.com/alextanhongpin/go-microservice/presentation/api"
	"github.com/alextanhongpin/pkg/gojwt"
	"github.com/pkg/errors"
)

// Service represents the service layer that handles business logic that is not
// contained in the domain entity.
type Service struct {
	signer gojwt.Signer
}

// NewService returns a new Service.
func NewService(signer gojwt.Signer) *Service {
	return &Service{signer}
}

// CreateAccessToken returns a new access token with the default claims given a
// user id.
func (s *Service) CreateAccessToken(userID string) (string, error) {
	if gostrings.IsEmpty(userID) {
		return "", errors.New("user_id is required")
	}
	accessToken, err := s.signer.Sign(func(c *gojwt.Claims) error {
		c.StandardClaims.Subject = userID
		c.Scope = api.Scopes(api.ScopeProfile, api.ScopeOpenID)
		// NOTE: Setting roles here can be problematic, especially when
		// the roles are dynamic.
		// TODO: Determine role based on user role.
		c.Role = api.RoleUser.String()
		return nil
	})
	return accessToken, errors.Wrap(err, "sign token failed")
}
