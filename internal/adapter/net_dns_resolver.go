package adapter

import (
	"context"
	"net"

	"github.com/GoSeoTaxi/email-validator/internal/domain"
)

type NetDNSResolver struct {
	resolver *net.Resolver
}

func NewNetDNSResolver(resolver *net.Resolver) domain.DNSResolver {
	return &NetDNSResolver{resolver: resolver}
}

func (n *NetDNSResolver) LookupMX(ctx context.Context, name string) ([]*net.MX, error) {
	return n.resolver.LookupMX(ctx, name)
}
