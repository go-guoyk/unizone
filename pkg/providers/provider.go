package providers

import (
	"context"
	"errors"
	"go.guoyk.net/sugar"
	"sync"
)

const (
	UnizoneNameKey = "unizone-name"
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
	Region      string
	TokenID     string
	TokenSecret string
	Logger      sugar.Logger
}

type Service struct {
	Name    string
	IP      string
	Comment string
}

type Provider interface {
	ListServices(ctx context.Context, network string, service string) ([]Service, error)
}

type Registration interface {
	CreateProvider(opts Options) (Provider, error)
}
