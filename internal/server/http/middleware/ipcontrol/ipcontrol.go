package ipcontrol

import (
	"net"
	"net/http"
)

type IPControl struct {
	TrustedSubnet *net.IPNet
}

func New(ts *net.IPNet) IPControl {
	return IPControl{
		TrustedSubnet: ts,
	}
}

func (mdl IPControl) AllowOnlyTrustedSubnet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if mdl.TrustedSubnet == nil {
			next.ServeHTTP(response, request)

			return
		}

		agentIP := mdl.getAgentIP(request)
		if agentIP == nil {
			response.WriteHeader(http.StatusForbidden)
			return
		}

		if mdl.TrustedSubnet.Contains(agentIP) {
			next.ServeHTTP(response, request)
			return
		}

		response.WriteHeader(http.StatusForbidden)
	})
}

func (mdl IPControl) getAgentIP(request *http.Request) net.IP {
	ip := request.Header.Get("X-Real-IP")
	return net.ParseIP(ip)
}
