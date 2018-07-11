package auth

import (
	"log"
	"net/http"
	"os"

	oidc "github.com/coreos/go-oidc"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	
	httperror "github.com/portainer/portainer/http/error"
)

var (
	clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
)

// GET request on /api/oauth/auth
func (handler *Handler) oauthAuthenticate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		log.Fatal(err)
	}
	// oidcConfig := &oidc.Config{
		// ClientID: clientID,
	// }
	// verifier := provider.Verifier(oidcConfig)

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://127.0.0.1:9000/auth/oauth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	state := "foobar" // Don't do this in production.
	
	http.Redirect(w, r, config.AuthCodeURL(state), http.StatusFound)
	
	
	return nil
}
