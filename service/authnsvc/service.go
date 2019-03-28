// Attempted the following name for package:
// - authenticator: this sounds more like a verb
// - authentication: too long
// - userlogin: is too specific, since user can also register
// - loginUser: breaks the convention, since package name is preferable a noun.
// - authz and authn is better.

package authnsvc

import "github.com/alextanhongpin/pkg/gojwt"

type Service struct {
	Login    LoginUseCase
	Register RegisterUseCase
}

func NewService(repo Repository, signer gojwt.Signer) *Service {
	createAccessToken := NewCreateAccessTokenUseCase(signer)
	return &Service{
		Login:    NewLoginUseCase(repo, createAccessToken),
		Register: NewRegisterUseCase(repo, createAccessToken),
	}
}
