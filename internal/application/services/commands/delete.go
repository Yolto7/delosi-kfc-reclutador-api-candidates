package commands

import (
	"context"

	"github.com/Yolto7/api-candidates/internal/domain/config"
	"github.com/Yolto7/api-candidates/internal/domain/repositories"
	"github.com/Yolto7/api-candidates/pkg/domain/logger"
)

// =====================================================================
// DTOs and Input/Output types
// =====================================================================

// DeleteServiceInput represents the input for candidate delete operation
type DeleteServiceInput struct {
	ID 						string `json:"id" validate:"required,notblank"`
}

// DeleteServiceOutput represents the output of candidate delete operation
type DeleteServiceOutput struct {
}

// =====================================================================
// Service Configuration
// =====================================================================

// DeleteService handles candidate delete operations with optimized performance
type DeleteService struct {
	config                 *config.Config
	logger                 logger.Logger
	candidateRepository repositories.CandidateRepository
}

// DeleteServiceConfig holds the configuration dependencies for DeleteService
type DeleteServiceConfig struct {
	Config                 *config.Config
	Logger                 logger.Logger
	CandidateRepository repositories.CandidateRepository
}

// NewDeleteService deletes a new instance of DeleteService with provided configuration
func NewDeleteService(cfg DeleteServiceConfig) *DeleteService {
	return &DeleteService{
		config:                 cfg.Config,
		logger:                 cfg.Logger,
		candidateRepository: cfg.CandidateRepository,
	}
}

// =====================================================================
// Main Service Logic
// =====================================================================

// Execute performs an optimized delete operation on candidate
// Returns error if validation fails or repository operations fail
func (svc *DeleteService) Execute(ctx context.Context, input *DeleteServiceInput) (*DeleteServiceOutput, error) {
	if err := svc.candidateRepository.Delete(ctx, input.ID); err != nil {
		return nil, err
	}

	return &DeleteServiceOutput{}, nil
}

