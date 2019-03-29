// Attempted the following name for package:
// - authenticator: this sounds more like a verb
// - authentication: too long
// - userlogin: is too specific, since user can also register
// - loginUser: breaks the convention, since package name is preferable a noun.
// - authz and authn is better.

package authnsvc

import (
	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/govalidator"
	"github.com/alextanhongpin/pkg/gojwt"
)

type (
	// CreateAccessTokenUseCase creates a new access token that last for
	// the given duration for a specific user.
	CreateAccessTokenUseCase struct {
		signer gojwt.Signer
	}
)

// NewCreateAccessTokenUseCase returns a new usecase to create access token.
func NewCreateAccessTokenUseCase(signer gojwt.Signer) *CreateAccessTokenUseCase {
	return &CreateAccessTokenUseCase{signer}
}

// CreateAccessToken creates a new token for the given user.
func (c *CreateAccessTokenUseCase) CreateAccessToken(userID string) (string, error) {
	if err := govalidator.Validate.Var(userID, "required"); err != nil {
		return "", err
	}
	accessToken, err := c.signer.Sign(func(c *gojwt.Claims) error {
		c.StandardClaims.Subject = userID
		c.Scope = api.Scopes(api.ScopeProfile, api.ScopeOpenID)
		// TODO: Determine role based on user role.
		c.Role = api.RoleUser.String()
		return nil
	})
	return accessToken, errors.Wrap(err, "sign token failed")
}
