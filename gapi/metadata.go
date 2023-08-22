package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader     = "user-agent"
	xForwarderForHeader = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (s *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if forwardedFor := md.Get(xForwarderForHeader); len(forwardedFor) > 0 {
			mtdt.ClientIp = forwardedFor[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIp = p.Addr.String()
	}

	return mtdt
}
