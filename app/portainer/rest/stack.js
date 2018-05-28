angular.module('portainer.app')
.factory('Stack', ['$resource', 'EndpointProvider', 'API_ENDPOINT_STACKS', function StackFactory($resource, EndpointProvider, API_ENDPOINT_STACKS) {
  'use strict';
  return $resource(API_ENDPOINT_STACKS + '/:id/:action', {},
  {
    get: { method: 'GET', params: { id: '@id' } },
    query: { method: 'GET', isArray: true },
    create: { method: 'POST', ignoreLoadingBar: true },
    update: { method: 'PUT', params: { id: '@id' }, ignoreLoadingBar: true },
    remove: { method: 'DELETE', params: { id: '@id'} },
    getStackFile: { method: 'GET', params: { id : '@id', action: 'file' } }
  });
}]);
