package auth

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/portainer/portainer"
	httperror "github.com/portainer/portainer/http/error"
	"github.com/portainer/portainer/http/request"
	"github.com/portainer/portainer/http/response"
)

type authenticatePayload struct {
	Username string
	Password string
}

type authenticateResponse struct {
	JWT string `json:"jwt"`
}

func (payload *authenticatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Username) {
		return portainer.Error("Invalid username")
	}
	if govalidator.IsNull(payload.Password) {
		return portainer.Error("Invalid password")
	}
	return nil
}

func (handler *Handler) authenticate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	if handler.authDisabled {
		return &httperror.HandlerError{http.StatusServiceUnavailable, "Cannot authenticate user. Portainer was started with the --no-auth flag", ErrAuthDisabled}
	}

	var payload authenticatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	settings, err := handler.SettingsService.Settings()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve settings from the database", err}
	}

	if settings.AuthenticationMethod == portainer.AuthenticationLDAP {
		err = handler.LDAPService.AuthenticateUser(payload.Username, payload.Password, &settings.LDAPSettings)
		if err != nil && err != portainer.ErrInvalidUsername {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to authenticate user via LDAP/AD", err}
		}

		u, err := handler.UserService.UserByUsername(payload.Username)
		if err != nil && u == nil {
			if err == portainer.ErrObjectNotFound {
				user := &portainer.User{
					Username: payload.Username,
					Role:     portainer.StandardUserRole,
				}

				err = handler.UserService.CreateUser(user)
				if err != nil {
					return &httperror.HandlerError{http.StatusInternalServerError, "Unable to create user", err}
				}
				
				u, err = handler.UserService.UserByUsername(payload.Username)
			} else {
				return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve a user with the specified username from the database", err}
			}
		}
		
		user := &portainer.User{
			Username: u.Username,
			Role:     portainer.StandardUserRole,
		}
		
		err = handler.AddUserIntoTeams(user, settings)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to add user to team", err}
		}
	}

	u, err := handler.UserService.UserByUsername(payload.Username)
	if err == portainer.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid credentials", ErrInvalidCredentials}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve a user with the specified username from the database", err}
	}

	if settings.AuthenticationMethod == portainer.AuthenticationInternal {
		err = handler.CryptoService.CompareHashAndData(u.Password, payload.Password)
		if err != nil {
			return &httperror.HandlerError{http.StatusUnprocessableEntity, "Invalid credentials", ErrInvalidCredentials}
		}
	}

	tokenData := &portainer.TokenData{
		ID:       u.ID,
		Username: u.Username,
		Role:     u.Role,
	}

	token, err := handler.JWTService.GenerateToken(tokenData)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to generate JWT token", err}
	}

	return response.JSON(w, &authenticateResponse{JWT: token})
}

func (handler *Handler) AddUserIntoTeams(user *portainer.User, settings *portainer.Settings) error {
	teams, err := handler.TeamService.Teams()
	if err != nil {
		return err
	}

	var userGroups []string
	if settings.AuthenticationMethod == portainer.AuthenticationLDAP {
		userGroups, err = handler.LDAPService.GetUserGroups(user.Username, &settings.LDAPSettings)
		if err != nil {
			return err
		}
	} else {
		return portainer.Error("Function AddUserIntoTeams on supports LDAP authentication")
	}

	userMemberships, err := handler.TeamMembershipService.TeamMembershipsByUserID(user.ID)
	if err != nil {
		return err
	}

	for _, team := range teams {
		if teamExists(team.Name, userGroups) {

			if teamMembershipExists(team.ID, userMemberships) {
				continue
			}

			membership := &portainer.TeamMembership{
				UserID: user.ID,
				TeamID: team.ID,
				Role:   portainer.TeamMember,
			}

			handler.TeamMembershipService.CreateTeamMembership(membership)
		}
	}
	return nil
}

func teamExists(teamName string, ldapGroups []string) bool {
	for _, group := range ldapGroups {
		if group == teamName {
			return true
		}
	}
	return false
}

func teamMembershipExists(teamID portainer.TeamID, memberships []portainer.TeamMembership) bool {
	for _, membership := range memberships {
		if membership.TeamID == teamID {
			return true
		}
	}
	return false
}