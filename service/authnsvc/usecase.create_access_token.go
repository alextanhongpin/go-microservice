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
	"github.com/alextanhongpin/pkg/gojwt"
)

type (
	createAccessTokenUseCase interface {
		CreateAccessToken(user string) (token string, err error)
	}
	CreateAccessTokenUseCase struct {
		signer gojwt.Signer
	}
)

func NewCreateAccessTokenUseCase(signer gojwt.Signer) *CreateAccessTokenUseCase {
	return &CreateAccessTokenUseCase{signer}
}

func (c *CreateAccessTokenUseCase) CreateAccessToken(user string) (string, error) {
	if len(user) == 0 {
		return "", errors.New("user is required")
	}
	accessToken, err := c.signer.Sign(func(c *gojwt.Claims) error {
		c.StandardClaims.Subject = user
		c.Scope = api.Scopes(api.ScopeProfile, api.ScopeOpenID)
		// TODO: Determine role based on user role.
		c.Role = api.RoleUser.String()
		return nil
	})
	return accessToken, errors.Wrap(err, "sign token failed")
}
