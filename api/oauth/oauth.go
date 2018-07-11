package oauth

import (
	"errors"
	"github.com/portainer/portainer"
)

// Service represents a service used to authenticate users against a OAuth.
type Service struct{}

// AuthenticateUser is used to authenticate a user against a OAuth.
func (*Service) AuthenticateUser(username string, settings *portainer.OAuthSettings) error {

	err := errors.New("oauth.go - AuthenticateUser is not yet implemented")
	return err
}

// TestConnectivity is used to test a connection against the OAuth server using the credentials
// specified in the OAuthSettings.
func (*Service) TestConnectivity(settings *portainer.OAuthSettings) error {
	err := errors.New("oauth.go - TestConnectivity is not yet implemented")
	return err
}
