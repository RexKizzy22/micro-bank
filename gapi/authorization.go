package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/Rexkizzy22/micro-bank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader     = "authorization"
	authorizationBearerType = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if values == nil {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearerType {
		return nil, fmt.Errorf("unsupported authorization type %s", authType)
	}

	bearerToken := fields[1]
	payload, err := server.token.VerifyToken(bearerToken)
	if err != nil {
		return nil, fmt.Errorf("unable to verify token")
	}

	return payload, nil
}
