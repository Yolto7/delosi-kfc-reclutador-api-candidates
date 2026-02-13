package repositories

import (
	"context"

	"github.com/Yolto7/api-candidates/internal/domain/entities"
)

type CandidateRepository interface {
	GetByID(ctx context.Context, id string) (*entities.Candidate, error)
	GetByCompositeKey(ctx context.Context, compositeKey string) (*entities.Candidate, error)
	Create(ctx context.Context, candidate *entities.Candidate) error
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}
