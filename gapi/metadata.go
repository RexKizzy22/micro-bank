package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) getMetadata(ctx context.Context) *Metadata {
	metad := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			metad.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			metad.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(xForwardedForHeader); len(userAgents) > 0 {
			metad.ClientIP = userAgents[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		metad.ClientIP = p.Addr.String()
	}

	return metad
}
