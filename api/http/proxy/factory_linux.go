// +build linux

package proxy

import (
	"net/http"
)

func (factory *proxyFactory) newNamedPipeProxy(path string) http.Handler {
	proxy := &localProxy{}
	return proxy
}
