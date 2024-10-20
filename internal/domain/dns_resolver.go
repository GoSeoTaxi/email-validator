package domain

import (
	"context"
	"net"
)

type DNSResolver interface {
	LookupMX(ctx context.Context, name string) ([]*net.MX, error)
}
