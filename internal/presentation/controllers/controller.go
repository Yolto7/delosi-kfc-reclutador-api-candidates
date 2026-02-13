package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Yolto7/api-candidates/internal/application/services/commands"
	"github.com/Yolto7/api-candidates/internal/application/services/queries"
	"github.com/Yolto7/api-candidates/internal/presentation/validators"
	errorCustom "github.com/Yolto7/api-candidates/pkg/domain/error"
	"github.com/Yolto7/api-candidates/pkg/domain/logger"
	"github.com/Yolto7/api-candidates/pkg/infrastructure/response"
	"github.com/aws/aws-lambda-go/events"
)


type CandidateController struct {
	logger 					logger.Logger
	getByIDService 	*queries.GetByIDService
	createService 	*commands.CreateService
	deleteService 	*commands.DeleteService
}

type CandidateControllerConfig struct {
	Logger 		logger.Logger
	GetByIDService 	*queries.GetByIDService
	CreateService		*commands.CreateService
	DeleteService		*commands.DeleteService
}

func NewCandidateController(cfg CandidateControllerConfig) *CandidateController {
  return &CandidateController{
		logger: 					cfg.Logger,
		getByIDService: 	cfg.GetByIDService,
		createService: 		cfg.CreateService,
		deleteService: 		cfg.DeleteService,
  }
}

// Queries
func (ctr *CandidateController) GetByID(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	id, ok := event.PathParameters["id"]
	if !ok || id == "" {
		return nil, errorCustom.NewError(errorCustom.BAD_REQUEST, "Invalid ID", "ERR_INVALID_ID")
	}

	req := queries.GetByIDServiceInput{
		ID: id,
	}
	if err := validators.GetByID(req); err != nil {
		return nil, err
	}

	ctr.logger.Info(fmt.Sprintf("GetByID request: %+v", req))
	result, err := ctr.getByIDService.Execute(ctx, req)
	if err != nil {
		return nil, errorCustom.FromError(err) 
	}

	ctr.logger.Info(fmt.Sprintf("GetByID result: %+v", result))
	return response.Success(http.StatusOK, "Got candidate successfully", result)
}

// Commands
func (ctr *CandidateController) Create(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var req commands.CreateServiceInput
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return nil, errorCustom.NewError(errorCustom.BAD_REQUEST, "Invalid JSON format", "ERR_INVALID_JSON")
	}
	if err := validators.Create(&req); err != nil {
		return nil, err
	}

	ctr.logger.Info(fmt.Sprintf("Create request: %+v", req))
	result, err := ctr.createService.Execute(ctx, &req)
	if err != nil {
		return nil, errorCustom.FromError(err) 
	}

	ctr.logger.Info(fmt.Sprintf("Create result: %+v", result))
	return response.Success(http.StatusOK, "Create candidate successfully", result)
}

func (ctr *CandidateController) Delete(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	id, ok := event.PathParameters["id"]
	if !ok || id == "" {
		return nil, errorCustom.NewError(errorCustom.BAD_REQUEST, "Invalid ID", "ERR_INVALID_ID")
	}

	req := commands.DeleteServiceInput{
		ID: id,
	}
	if err := validators.Delete(&req); err != nil {
		return nil, err
	}

	ctr.logger.Info(fmt.Sprintf("Delete request: %+v", req))
	result, err := ctr.deleteService.Execute(ctx, &req)
	if err != nil {
		return nil, errorCustom.FromError(err) 
	}

	ctr.logger.Info(fmt.Sprintf("Delete result: %+v", result))
	return response.Success(http.StatusOK, "Delete candidate successfully", result)
}	

