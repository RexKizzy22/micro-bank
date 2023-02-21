package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	gRPCGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

// gets metadata from HTTP and gRPC requests to a gRPC-gateway server using the metadata and peer packages
func (server *Server) getMetadata(ctx context.Context) *Metadata {
	metad := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// extract metadata (user agent) from a HTTP request to a gRPC-gateway server
		if userAgents := md.Get(gRPCGatewayUserAgentHeader); len(userAgents) > 0 {
			metad.UserAgent = userAgents[0]
		}

		// extract metadata (user agent) from a gRPC request to a gRPC-gateway server
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			metad.UserAgent = userAgents[0]
		}

		// extract metadata (client IP) from a HTTP request to a gRPC-gateway server
		if userAgents := md.Get(xForwardedForHeader); len(userAgents) > 0 {
			metad.ClientIP = userAgents[0]
		}
	}

	// extract metadata (client IP address) from a gRPC request to a gRPC-gateway server
	if p, ok := peer.FromContext(ctx); ok {
		metad.ClientIP = p.Addr.String()
	}

	return metad
}
