package aws

import (
	"context"

	"github.com/infrasphere/control-plane/backend/internal/providers/mock"
	"github.com/infrasphere/control-plane/backend/pkg/cloud"
)

type Provider struct {
	*mock.Provider
}

func New() *Provider {
	return &Provider{Provider: mock.New()}
}

func (p *Provider) Name() string { return "aws" }

func (p *Provider) Connect(ctx context.Context, cfg cloud.ProviderConfig) error {
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	return p.Provider.Connect(ctx, cfg)
}

