// +build windows

package proxy

import (
	"net"
	"net/http"

	"github.com/Microsoft/go-winio"
)

func (factory *proxyFactory) newNamedPipeProxy(path string) http.Handler {
	proxy := &localProxy{}
	transport := &proxyTransport{
		ResourceControlService: factory.ResourceControlService,
		TeamMembershipService:  factory.TeamMembershipService,
		SettingsService:        factory.SettingsService,
		dockerTransport:        newNamedPipeTransport(path),
	}
	proxy.Transport = transport
	return proxy
}

func newNamedPipeTransport(namedPipePath string) *http.Transport {
	return &http.Transport{
		Dial: func(proto, addr string) (conn net.Conn, err error) {
			return winio.DialPipe(namedPipePath, nil)
		},
	}
}
