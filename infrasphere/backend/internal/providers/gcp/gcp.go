package gcp

import "github.com/infrasphere/control-plane/backend/internal/providers/mock"

type Provider struct{ *mock.Provider }

func New() *Provider { return &Provider{Provider: mock.New()} }
func (p *Provider) Name() string { return "gcp" }

