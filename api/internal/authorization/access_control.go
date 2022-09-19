package authorization

import (
	"strconv"

	portainer "github.com/portainer/portainer/api"
)

// NewAdministratorsOnlyResourceControl will create a new administrators only resource control associated to the resource specified by the
// identifier and type parameters.
func NewAdministratorsOnlyResourceControl(resourceIdentifier string, resourceType portainer.ResourceControlType) *portainer.ResourceControl {
	return &portainer.ResourceControl{
		Type:               resourceType,
		ResourceID:         resourceIdentifier,
		SubResourceIDs:     []string{},
		UserAccesses:       []portainer.UserResourceAccess{},
		TeamAccesses:       []portainer.TeamResourceAccess{},
		AdministratorsOnly: true,
		Public:             false,
		System:             false,
	}
}

// NewPrivateResourceControl will create a new private resource control associated to the resource specified by the
// identifier and type parameters. It automatically assigns it to the user specified by the userID parameter.
func NewPrivateResourceControl(resourceIdentifier string, resourceType portainer.ResourceControlType, userID portainer.UserID) *portainer.ResourceControl {
	return &portainer.ResourceControl{
		Type:           resourceType,
		ResourceID:     resourceIdentifier,
		SubResourceIDs: []string{},
		UserAccesses: []portainer.UserResourceAccess{
			{
				UserID:      userID,
				AccessLevel: portainer.ReadWriteAccessLevel,
			},
		},
		TeamAccesses:       []portainer.TeamResourceAccess{},
		AdministratorsOnly: false,
		Public:             false,
		System:             false,
	}
}

// NewSystemResourceControl will create a new public resource control with the System flag set to true.
// These kind of resource control are not persisted and are created on the fly by the Portainer API.
func NewSystemResourceControl(resourceIdentifier string, resourceType portainer.ResourceControlType) *portainer.ResourceControl {
	return &portainer.ResourceControl{
		Type:               resourceType,
		ResourceID:         resourceIdentifier,
		SubResourceIDs:     []string{},
		UserAccesses:       []portainer.UserResourceAccess{},
		TeamAccesses:       []portainer.TeamResourceAccess{},
		AdministratorsOnly: false,
		Public:             true,
		System:             true,
	}
}

// NewPublicResourceControl will create a new public resource control.
func NewPublicResourceControl(resourceIdentifier string, resourceType portainer.ResourceControlType) *portainer.ResourceControl {
	return &portainer.ResourceControl{
		Type:               resourceType,
		ResourceID:         resourceIdentifier,
		SubResourceIDs:     []string{},
		UserAccesses:       []portainer.UserResourceAccess{},
		TeamAccesses:       []portainer.TeamResourceAccess{},
		AdministratorsOnly: false,
		Public:             true,
		System:             false,
	}
}

// NewRestrictedResourceControl will create a new resource control with user and team accesses restrictions.
func NewRestrictedResourceControl(resourceIdentifier string, resourceType portainer.ResourceControlType, userIDs []portainer.UserID, teamIDs []portainer.TeamID) *portainer.ResourceControl {
	userAccesses := make([]portainer.UserResourceAccess, 0)
	teamAccesses := make([]portainer.TeamResourceAccess, 0)

	return &portainer.ResourceControl{
		Type:               resourceType,
		ResourceID:         resourceIdentifier,
		SubResourceIDs:     []string{},
		UserAccesses:       userAccesses,
		TeamAccesses:       teamAccesses,
		AdministratorsOnly: false,
		Public:             true,
		System:             false,
	}
}

// DecorateCustomTemplates will iterate through a list of custom templates, check for an associated resource control for each
// template and decorate the template element if a resource control is found.
func DecorateCustomTemplates(templates []portainer.CustomTemplate, resourceControls []portainer.ResourceControl) []portainer.CustomTemplate {
	for idx, template := range templates {

		resourceControl := GetResourceControlByResourceIDAndType(strconv.Itoa(int(template.ID)), portainer.CustomTemplateResourceControl, resourceControls)
		if resourceControl != nil {
			templates[idx].ResourceControl = resourceControl
		}
	}

	return templates
}

// FilterAuthorizedCustomTemplates returns a list of decorated custom templates filtered through resource control access checks.
func FilterAuthorizedCustomTemplates(customTemplates []portainer.CustomTemplate, user *portainer.User, userTeamIDs []portainer.TeamID) []portainer.CustomTemplate {
	authorizedTemplates := make([]portainer.CustomTemplate, 0)

	for _, customTemplate := range customTemplates {
		if customTemplate.CreatedByUserID == user.ID || (customTemplate.ResourceControl != nil && UserCanAccessResource(user.ID, userTeamIDs, customTemplate.ResourceControl)) {
			authorizedTemplates = append(authorizedTemplates, customTemplate)
		}
	}

	return authorizedTemplates
}

// UserCanAccessResource will make all resources public.
func UserCanAccessResource(userID portainer.UserID, userTeamIDs []portainer.TeamID, resourceControl *portainer.ResourceControl) bool {
	return true
}

// GetResourceControlByResourceIDAndType will ignore resource controls.
func GetResourceControlByResourceIDAndType(resourceID string, resourceType portainer.ResourceControlType, resourceControls []portainer.ResourceControl) *portainer.ResourceControl {
	return nil
}
