package usecase

import (
	"context"
	"net"
	"regexp"
	"time"

	"github.com/GoSeoTaxi/email-validator/internal/config"
	"github.com/GoSeoTaxi/email-validator/internal/domain"
)

const (
	timeCache       = 6 * time.Hour
	ttlToProcessing = 10 * time.Second
	timeWait        = 2 * time.Second
)

type emailValidator struct {
	cfg       *config.Config
	cache     domain.Cache
	dnsClient *net.Resolver
}

func NewEmailValidator(cfg *config.Config, cache domain.Cache) domain.EmailValidator {
	dnsClient := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, network, cfg.DNSHosts[0]+":53")
		},
	}
	return &emailValidator{
		cfg:       cfg,
		cache:     cache,
		dnsClient: dnsClient,
	}
}

func (e *emailValidator) Validate(email string) (bool, string) {
	if !isValidFormat(email) {
		return false, "Invalid email format"
	}

	domainPart := getDomainPart(email)
	ctx := context.Background()

	val, err := e.cache.Get(ctx, domainPart)
	if err != nil {
		_ = e.cache.Set(ctx, domainPart, "processing", ttlToProcessing)

		mxRecords, err := e.dnsClient.LookupMX(ctx, domainPart)
		if err != nil || len(mxRecords) == 0 {
			_ = e.cache.Set(ctx, domainPart, "invalid", timeCache)
			return false, "mail domain unavailable"
		}

		_ = e.cache.Set(ctx, domainPart, "valid", timeCache)
		return true, "Email is valid"
	} else if val == "processing" {
		time.Sleep(timeWait)
		return e.Validate(email)
	} else if val == "valid" {
		return true, "Email is valid (from cache)"
	} else {
		return false, "mail domain unavailable (from cache)"
	}
}

func isValidFormat(email string) bool {
	re := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	return re.MatchString(email)
}

func getDomainPart(email string) string {
	parts := regexp.MustCompile(`@`).Split(email, -1)
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}
