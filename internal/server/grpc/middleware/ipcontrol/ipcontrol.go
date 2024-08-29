package ipcontrol

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

type Middleware struct {
	TrustedSubnet *net.IPNet
}

func New(ts *net.IPNet) *Middleware {
	return &Middleware{
		TrustedSubnet: ts,
	}
}

func (mdl *Middleware) IPControl(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if mdl.TrustedSubnet == nil {
		return handler(ctx, req)
	}

	agentIP := mdl.getAgentIP(ctx)
	if agentIP == nil {
		return nil, fmt.Errorf("server is configured with trusted subnet, but agent did not pass ip")
	}

	if mdl.TrustedSubnet.Contains(agentIP) {
		return handler(ctx, req)
	}

	return nil, fmt.Errorf("server is configured with trusted subnet that does not contain agent ip")
}

func (mdl *Middleware) SetTrustedSubnet(ts *net.IPNet) {
	mdl.TrustedSubnet = ts
}

func (mdl *Middleware) getAgentIP(ctx context.Context) net.IP {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	vals := md.Get("X-Real-IP")
	if len(vals) != 1 {
		return nil
	}
	return net.ParseIP(vals[0])
}
