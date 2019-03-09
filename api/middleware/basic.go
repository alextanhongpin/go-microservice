package middleware

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/pkg/set"
)

// BasicAuthorizer middleware represents the endpoint protected with the Authorization
// Header set to `Basic username:password`.
func BasicAuthorizer(credential string, extras ...string) gin.HandlerFunc {
	credManager := set.New()
	credentials := append([]string{credential}, extras...)
	for _, cred := range credentials {
		credManager.Add(cred)
	}
	checkCredential := func(auth string) error {
		paths := strings.Split(auth, " ")
		if len(paths) != 2 {
			return errors.New("missing authorization header")
		}
		bearer, token := paths[0], paths[1]
		if valid := strings.EqualFold(bearer, AuthorizationBasic); !valid {
			return errors.New("invalid authorization header")
		}
		// Decode the base64 encoded 'username:password'.
		b, err := base64.URLEncoding.DecodeString(token)
		if err != nil {
			return errors.Wrap(err, "decode token failed")
		}

		h := sha256.New()
		h.Write(b)
		hashed := hex.EncodeToString(h.Sum(nil))

		if !credManager.Has(hashed) {
			return errors.New("invalid credentials")
		}
		return nil
	}
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if err := checkCredential(auth); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api.NewError(c, err),
			)
			return
		}
		c.Next()
	}
}
