package commands

import (
	"context"
	"fmt"

	"github.com/Yolto7/api-candidates/internal/domain/config"
	"github.com/Yolto7/api-candidates/internal/domain/entities"
	"github.com/Yolto7/api-candidates/internal/domain/repositories"
	"github.com/Yolto7/api-candidates/pkg/domain/constants"
	errorCustom "github.com/Yolto7/api-candidates/pkg/domain/error"
	"github.com/Yolto7/api-candidates/pkg/domain/logger"
	"github.com/Yolto7/api-candidates/pkg/infrastructure/utils"
	pkgIUtils "github.com/Yolto7/api-candidates/pkg/infrastructure/utils"
)

// =====================================================================
// DTOs and Input/Output types
// =====================================================================

// CreateServiceInput represents the input for candidate create operation
type CreateServiceInput struct {
	ID                                string `json:"id" validate:"required,notblank"`
	SheetID                           string `json:"sheetId" validate:"required,notblank"`
	RowID                             string `json:"rowId" validate:"required,notblank"`
	ColumnPostulantResponseId         string `json:"columnPostulantResponseId" validate:"required,notblank"`
	ColumnPostulantDateTimeResponseId string `json:"columnPostulantDateTimeResponseId" validate:"required,notblank"`
	ColumnPostulantConfirmedId        string `json:"columnPostulantConfirmedId" validate:"required,notblank"`
}

// CreateServiceOutput represents the output of candidate create operation
type CreateServiceOutput struct {
}

// =====================================================================
// Service Configuration
// =====================================================================

// CreateService handles candidate create operations with optimized performance
type CreateService struct {
	config              *config.Config
	logger              logger.Logger
	candidateRepository repositories.CandidateRepository
}

// CreateServiceConfig holds the configuration dependencies for CreateService
type CreateServiceConfig struct {
	Config              *config.Config
	Logger              logger.Logger
	CandidateRepository repositories.CandidateRepository
}

// NewCreateService creates a new instance of CreateService with provided configuration
func NewCreateService(cfg CreateServiceConfig) *CreateService {
	return &CreateService{
		config:              cfg.Config,
		logger:              cfg.Logger,
		candidateRepository: cfg.CandidateRepository,
	}
}

// =====================================================================
// Main Service Logic
// =====================================================================

// Execute performs an optimized create operation on candidate
// Returns error if validation fails or repository operations fail
func (svc *CreateService) Execute(ctx context.Context, input *CreateServiceInput) (*CreateServiceOutput, error) {
	compositeKey := fmt.Sprintf("%s#%s", input.ID, input.SheetID)

	exists, err := svc.candidateRepository.GetByCompositeKey(ctx, compositeKey)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		svc.logger.Error(utils.NewSafeError(err, "Candidate with ID already exists"))
		return nil, errorCustom.NewError(errorCustom.BAD_REQUEST, "Candidate already exists", "ERR_CANDIDATE_EXISTS")
	}

	candidate := &entities.Candidate{
		ID:                                input.ID,
		CompositeKey:                      compositeKey,
		SheetID:                           input.SheetID,
		RowID:                             input.RowID,
		ColumnPostulantResponseId:         input.ColumnPostulantResponseId,
		ColumnPostulantDateTimeResponseId: input.ColumnPostulantDateTimeResponseId,
		ColumnPostulantConfirmedId:        input.ColumnPostulantConfirmedId,
		CreatedAt:                         pkgIUtils.NowDateTime("America/Lima"),
		CreatedBy:                         constants.SYSTEM_USER,
		Deleted:                           false,
	}

	if err := svc.candidateRepository.Create(ctx, candidate); err != nil {
		return nil, err
	}

	return &CreateServiceOutput{}, nil
}
