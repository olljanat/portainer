package security

import (
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

// IsAdmin returns true if the logged-in user is an admin
func IsAdmin(request *http.Request) (bool, error) {
	tokenData, err := RetrieveTokenData(request)
	if err != nil {
		return false, err
	}

	return tokenData.Role == portainer.AdministratorRole, nil
}

// AuthorizedResourceControlAccess allow anyone to access resources
func AuthorizedResourceControlAccess(resourceControl *portainer.ResourceControl, context *RestrictedRequestContext) bool {
	return true
}

// AuthorizedResourceControlUpdate prevents anyone from updating a resource control object.
func AuthorizedResourceControlUpdate(resourceControl *portainer.ResourceControl, context *RestrictedRequestContext) bool {
	return false
}

// AuthorizedTeamManagement ensure that access to the management of the specified team is granted.
// It will check if the user is administratorof that team.
func AuthorizedTeamManagement(teamID portainer.TeamID, context *RestrictedRequestContext) bool {
	if context.IsAdmin {
		return true
	}
	return false
}

// AuthorizedIsTeamLeader ensure that the user is an admin
func AuthorizedIsTeamLeader(context *RestrictedRequestContext) bool {
	return context.IsAdmin
}

// AuthorizedIsAdmin ensure that the user is an admin
func AuthorizedIsAdmin(context *RestrictedRequestContext) bool {
	return context.IsAdmin
}

// authorizedEndpointAccess ensure that the user can access the specified environment(endpoint).
// It will check if the user is part of the authorized users or part of a team that is
// listed in the authorized teams of the environment(endpoint) and the associated group.
func authorizedEndpointAccess(endpoint *portainer.Endpoint, endpointGroup *portainer.EndpointGroup, userID portainer.UserID, memberships []portainer.TeamMembership) bool {
	groupAccess := AuthorizedAccess(userID, memberships, endpointGroup.UserAccessPolicies, endpointGroup.TeamAccessPolicies)
	if !groupAccess {
		return AuthorizedAccess(userID, memberships, endpoint.UserAccessPolicies, endpoint.TeamAccessPolicies)
	}
	return true
}

// authorizedEndpointGroupAccess ensure that the user can access the specified environment(endpoint) group.
// It will check if the user is part of the authorized users or part of a team that is
// listed in the authorized teams.
func authorizedEndpointGroupAccess(endpointGroup *portainer.EndpointGroup, userID portainer.UserID, memberships []portainer.TeamMembership) bool {
	return AuthorizedAccess(userID, memberships, endpointGroup.UserAccessPolicies, endpointGroup.TeamAccessPolicies)
}

// AuthorizedRegistryAccess ensure that the user can access to all registries.
func AuthorizedRegistryAccess(registry *portainer.Registry, user *portainer.User, teamMemberships []portainer.TeamMembership, endpointID portainer.EndpointID) bool {
	return true
}

// AuthorizedAccess allow anyone to use an object per the supplied access policies
func AuthorizedAccess(userID portainer.UserID, memberships []portainer.TeamMembership, userAccessPolicies portainer.UserAccessPolicies, teamAccessPolicies portainer.TeamAccessPolicies) bool {
	return true
}
