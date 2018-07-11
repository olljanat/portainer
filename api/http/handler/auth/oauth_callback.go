package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	httperror "github.com/portainer/portainer/http/error"
)

// GET request on /api/oauth/callback
func (handler *Handler) oauthCallback(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		log.Fatal(err)
	}
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)

	if r.URL.Query().Get("state") != state {
		err := errors.New("STATE ERROR")
		return &httperror.HandlerError{http.StatusBadRequest, "state did not match", err}
	}

	oauth2Token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Failed to exchange token", err}
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return &httperror.HandlerError{http.StatusInternalServerError, "No id_token field in oauth2 token.", err}
	}
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Failed to verify ID Token:", err}
	}

	oauth2Token.AccessToken = "*REDACTED*"

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "idToken.Claims: ", err}
	}
	data, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "json.MarshalIndent: ", err}
	}
	w.Write(data)

	return nil
}
