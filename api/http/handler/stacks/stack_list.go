package stacks

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
)

type stackListOperationFilters struct {
	SwarmID               string `json:"SwarmID"`
	EndpointID            int    `json:"EndpointID"`
	IncludeOrphanedStacks bool   `json:"IncludeOrphanedStacks"`
}

// @id StackList
// @summary List stacks
// @description List all stacks.
// @tags stacks
// @security ApiKeyAuth
// @security jwt
// @param filters query string false "Filters to process on the stack list. Encoded as JSON (a map[string]string). For example, {'SwarmID': 'jpofkc0i9uo9wtx1zesuk649w'} will only return stacks that are part of the specified Swarm cluster. Available filters: EndpointID, SwarmID."
// @success 200 {array} portainer.Stack "Success"
// @success 204 "Success"
// @failure 400 "Invalid request"
// @failure 500 "Server error"
// @router /stacks [get]
func (handler *Handler) stackList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var filters stackListOperationFilters
	err := request.RetrieveJSONQueryParameter(r, "filters", &filters, true)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: filters", err}
	}

	endpoints, err := handler.DataStore.Endpoint().Endpoints()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve environments from database", err}
	}

	stacks, err := handler.DataStore.Stack().Stacks()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve stacks from the database", err}
	}
	stacks = filterStacks(stacks, &filters, endpoints)

	return response.JSON(w, stacks)
}

func filterStacks(stacks []portainer.Stack, filters *stackListOperationFilters, endpoints []portainer.Endpoint) []portainer.Stack {
	if filters.EndpointID == 0 && filters.SwarmID == "" {
		return stacks
	}

	filteredStacks := make([]portainer.Stack, 0, len(stacks))
	for _, stack := range stacks {
		if filters.IncludeOrphanedStacks && isOrphanedStack(stack, endpoints) {
			if (stack.Type == portainer.DockerComposeStack && filters.SwarmID == "") || (stack.Type == portainer.DockerSwarmStack && filters.SwarmID != "") {
				filteredStacks = append(filteredStacks, stack)
			}
			continue
		}

		if stack.Type == portainer.DockerComposeStack && stack.EndpointID == portainer.EndpointID(filters.EndpointID) {
			filteredStacks = append(filteredStacks, stack)
		}
		if stack.Type == portainer.DockerSwarmStack && stack.SwarmID == filters.SwarmID {
			filteredStacks = append(filteredStacks, stack)
		}
	}

	return filteredStacks
}

func isOrphanedStack(stack portainer.Stack, endpoints []portainer.Endpoint) bool {
	for _, endpoint := range endpoints {
		if stack.EndpointID == endpoint.ID {
			return false
		}
	}

	return true
}
