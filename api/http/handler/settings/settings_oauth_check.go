package settings

import (
	"net/http"

	"github.com/portainer/portainer"
	httperror "github.com/portainer/portainer/http/error"
	"github.com/portainer/portainer/http/request"
	"github.com/portainer/portainer/http/response"
)

type settingsOAuthCheckPayload struct {
	OAuthSettings portainer.OAuthSettings
}

func (payload *settingsOAuthCheckPayload) Validate(r *http.Request) error {
	return nil
}

// PUT request on /settings/oauth/check
func (handler *Handler) settingsOAuthCheck(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload settingsOAuthCheckPayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	err = handler.OAuthService.TestConnectivity(&payload.OAuthSettings)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to connect to OAuth server", err}
	}

	return response.Empty(w)
}
