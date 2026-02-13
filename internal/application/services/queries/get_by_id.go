package queries

import (
	"context"

	"github.com/Yolto7/api-candidates/internal/domain/repositories"
	errorCustom "github.com/Yolto7/api-candidates/pkg/domain/error"
)

type GetByIDServiceInput struct {
	ID string `json:"id" validate:"notblank"`
}

type GetByIDServiceOutput struct {
	ID           string `json:"id"`
	DocumentID   string `json:"documentId"`
	DocumentName string `json:"documentName"`
	SheetID      string `json:"sheetId"`
	SheetName    string `json:"sheetName"`
}

type GetByIDService struct {
	candidateRepository repositories.CandidateRepository
}

type GetByIDServiceConfig struct {
	CandidateRepository repositories.CandidateRepository
}

func NewGetByIDService(cfg GetByIDServiceConfig) *GetByIDService {
	return &GetByIDService{
		candidateRepository: cfg.CandidateRepository,
	}
}

func (svc *GetByIDService) Execute(ctx context.Context, input GetByIDServiceInput) (*GetByIDServiceOutput, error) {
	candidate, err := svc.candidateRepository.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if candidate == nil {
		return nil, errorCustom.NewError(errorCustom.BAD_REQUEST, "Candidate not found", "ERR_CANDIDATE_NOT_FOUND")
	}

	return &GetByIDServiceOutput{
		ID:           candidate.ID,
		DocumentID:   candidate.DocumentID,
		DocumentName: candidate.DocumentName,
		SheetID:      candidate.SheetID,
		SheetName:    candidate.SheetName,
	}, nil
}
