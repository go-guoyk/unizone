package providers

import (
	"context"
	"errors"
	"regexp"
	"sync"
)

const (
	UnizoneNameKey = "unizone-name"
)

var (
	regexpRecordName = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
)

var (
	ErrUnknownProvider = errors.New("unknown provider")
)

var (
	registrations     = map[string]Registration{}
	registrationsLock = &sync.RWMutex{}
)

func Register(name string, reg Registration) {
	registrationsLock.Lock()
	defer registrationsLock.Unlock()
	registrations[name] = reg
}

func Create(name string, opts Options) (Provider, error) {
	registrationsLock.RLock()
	defer registrationsLock.RUnlock()

	reg := registrations[name]
	if reg == nil {
		return nil, ErrUnknownProvider
	}
	return reg.CreateProvider(opts)
}

type Options struct {
	Name        string
	Region      string
	TokenID     string
	TokenSecret string
}

type Record struct {
	Name    string
	IP      string
	Comment string
}

type Provider interface {
	ListRecords(ctx context.Context, network string, service string) ([]Record, error)
}

type Registration interface {
	CreateProvider(opts Options) (Provider, error)
}
