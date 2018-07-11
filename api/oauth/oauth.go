package oauth

import (
	"errors"
// 	"strings"

// 	"github.com/coreos/go-oidc"
	"github.com/portainer/portainer"
)

// Service represents a service used to authenticate users against a OAuth.
type Service struct{}

// AuthenticateUser is used to authenticate a user against a OAuth.
func (*Service) AuthenticateUser(username string, settings *portainer.OAuthSettings) error {

	err := errors.New("oauth.go - OAuth authentication is not yet implemented")
	return err
}
