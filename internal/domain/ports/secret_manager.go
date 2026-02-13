package ports

import (
	"context"
	"time"
)

// SecretValue represents a secret with its metadata
type SecretValue struct {
	Value       string
	Version     string
	CreatedAt   time.Time
	LastUpdated time.Time
	ARN         string
	Tags        map[string]string
}

// SecretOptions provides configuration for secret operations
type SecretOptions struct {
	VersionID     string
	VersionStage  string
	ForceRefresh  bool
	Timeout       time.Duration
}

// SecretsManager defines the interface for secret management operations
type SecretsManager interface {
	GetSecret(ctx context.Context, secretName string, opts *SecretOptions) (*SecretValue, error)
}

type SecretCache interface {
	Get(key string) (*SecretValue, bool)
	Set(key string, value *SecretValue, ttl time.Duration)
	Delete(key string)
	Clear()
}
