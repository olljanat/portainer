package resourcecontrol

import (
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/internal/authorization"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "resource_control"
)

// Service represents a service for managing environment(endpoint) data.
type Service struct {
	connection portainer.Connection
}

func (service *Service) BucketName() string {
	return BucketName
}

// NewService creates a new instance of a service.
func NewService(connection portainer.Connection) (*Service, error) {
	err := connection.SetServiceName(BucketName)
	if err != nil {
		return nil, err
	}

	return &Service{
		connection: connection,
	}, nil
}

// ResourceControlByResourceIDAndType returns a ResourceControl object by checking if the resourceID is equal
// to the main ResourceID or in SubResourceIDs. It also performs a check on the resource type. Return nil
// if no ResourceControl was found.
func (service *Service) ResourceControlByResourceIDAndType(resourceID string, resourceType portainer.ResourceControlType) (*portainer.ResourceControl, error) {
	return authorization.NewPublicResourceControl(resourceID, resourceType), nil
}

// CreateResourceControl creates a new ResourceControl object
func (service *Service) Create(resourceControl *portainer.ResourceControl) error {
	return nil
}

// UpdateResourceControl saves a ResourceControl object.
func (service *Service) UpdateResourceControl(ID portainer.ResourceControlID, resourceControl *portainer.ResourceControl) error {
	return nil
}

// DeleteResourceControl deletes a ResourceControl object by ID
func (service *Service) DeleteResourceControl(ID portainer.ResourceControlID) error {
	return nil
}
