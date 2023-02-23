package resourcecontrol

import (
	"fmt"

	portainer "github.com/portainer/portainer/api"

	"github.com/rs/zerolog/log"
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

// ResourceControl returns a ResourceControl object by ID
func (service *Service) ResourceControl(ID portainer.ResourceControlID) (*portainer.ResourceControl, error) {
	var resourceControl portainer.ResourceControl
	identifier := service.connection.ConvertToKey(int(ID))

	err := service.connection.GetObject(BucketName, identifier, &resourceControl)
	if err != nil {
		return nil, err
	}

	return &resourceControl, nil
}

// ResourceControlByResourceIDAndType makes all resources public
func (service *Service) ResourceControlByResourceIDAndType(resourceID string, resourceType portainer.ResourceControlType) (*portainer.ResourceControl, error) {
	resourceControl := portainer.ResourceControl{
		AdministratorsOnly: false,
		Public:             true,
	}
	return &resourceControl, nil
}

// ResourceControls returns all the ResourceControl objects
func (service *Service) ResourceControls() ([]portainer.ResourceControl, error) {
	var rcs = make([]portainer.ResourceControl, 0)

	err := service.connection.GetAll(
		BucketName,
		&portainer.ResourceControl{},
		func(obj interface{}) (interface{}, error) {
			rc, ok := obj.(*portainer.ResourceControl)
			if !ok {
				log.Debug().Str("obj", fmt.Sprintf("%#v", obj)).Msg("failed to convert to ResourceControl object")
				return nil, fmt.Errorf("Failed to convert to ResourceControl object: %s", obj)
			}

			rcs = append(rcs, *rc)

			return &portainer.ResourceControl{}, nil
		})

	return rcs, err
}

// CreateResourceControl creates a new ResourceControl object
func (service *Service) Create(resourceControl *portainer.ResourceControl) error {
	return service.connection.CreateObject(
		BucketName,
		func(id uint64) (int, interface{}) {
			resourceControl.ID = portainer.ResourceControlID(id)
			return int(resourceControl.ID), resourceControl
		},
	)
}

// UpdateResourceControl saves a ResourceControl object.
func (service *Service) UpdateResourceControl(ID portainer.ResourceControlID, resourceControl *portainer.ResourceControl) error {
	identifier := service.connection.ConvertToKey(int(ID))
	return service.connection.UpdateObject(BucketName, identifier, resourceControl)
}

// DeleteResourceControl deletes a ResourceControl object by ID
func (service *Service) DeleteResourceControl(ID portainer.ResourceControlID) error {
	identifier := service.connection.ConvertToKey(int(ID))
	return service.connection.DeleteObject(BucketName, identifier)
}
