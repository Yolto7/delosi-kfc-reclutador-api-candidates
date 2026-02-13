package container

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/Yolto7/api-candidates/internal/application/services/commands"
	"github.com/Yolto7/api-candidates/internal/application/services/queries"
	dConfig "github.com/Yolto7/api-candidates/internal/domain/config"
	iConfig "github.com/Yolto7/api-candidates/internal/infrastructure/config"
	iRepositories "github.com/Yolto7/api-candidates/internal/infrastructure/repositories"
	"github.com/Yolto7/api-candidates/internal/presentation/controllers"
	pkgDLogger "github.com/Yolto7/api-candidates/pkg/domain/logger"
	pkgILogger "github.com/Yolto7/api-candidates/pkg/infrastructure/logger"
	"github.com/Yolto7/api-candidates/pkg/infrastructure/middlewares"
	"github.com/Yolto7/api-candidates/pkg/infrastructure/persistence/dynamo"
)

// =============================================================================
// CONTAINERS OPTIMIZADOS - Solo resuelven dependencias específicas
// =============================================================================

// MainLambdaContainer - Solo dependencias mínimas para GET/POST
type MainLambdaContainer struct {
	ctx     context.Context
	logger  pkgDLogger.Logger
	config  *dConfig.Config
	
	// Lazy loading con sync.Once para thread safety
	dynamoOnce       sync.Once
	dynamoClient     *dynamodb.Client
	dynamoErr        error
	
	controllersOnce  sync.Once
	controller       *controllers.CandidateController
	controllersErr   error
	
	middlewaresOnce  sync.Once
	traceMiddleware  middlewares.Middleware
	baseMiddleware   middlewares.Middleware
	errorMiddleware  middlewares.Middleware
}

func NewMainLambdaContainer(ctx context.Context) (*MainLambdaContainer, error) {
	config, err := iConfig.Load()
	if err != nil {
		return nil, err
	}
	
	return &MainLambdaContainer{
		ctx:    ctx,
		logger: pkgILogger.NewZeroLogLogger(),
		config: config,
	}, nil
}

func (c *MainLambdaContainer) getDynamoClient() (*dynamodb.Client, error) {
	c.dynamoOnce.Do(func() {
		c.dynamoClient, c.dynamoErr = dynamo.GetClient(c.ctx)
	})
	return c.dynamoClient, c.dynamoErr
}

func (c *MainLambdaContainer) Logger() pkgDLogger.Logger {
	return c.logger
}

func (c *MainLambdaContainer) GetCandidateController() (*controllers.CandidateController, error) {
	c.controllersOnce.Do(func() {
		dynamoClient, err := c.getDynamoClient()
		if err != nil {
			c.controllersErr = err
			return
		}
		
		candidateRepo := iRepositories.NewCandidateDynamoRepository(c.logger, dynamoClient, c.config.CANDIDATES_TABLE_NAME)
		
		getByIDService := queries.NewGetByIDService(queries.GetByIDServiceConfig{
			CandidateRepository: candidateRepo,
		})
		
		createService := commands.NewCreateService(commands.CreateServiceConfig{
			Config:                 c.config,
			Logger:         				c.logger,
			CandidateRepository: candidateRepo,
		})

		deleteService := commands.NewDeleteService(commands.DeleteServiceConfig{
			Config:                 c.config,
			Logger:         				c.logger,
			CandidateRepository: candidateRepo,
		})
		
		c.controller = controllers.NewCandidateController(controllers.CandidateControllerConfig{
			Logger:         c.logger,
			GetByIDService: getByIDService,
			CreateService:  createService,
			DeleteService:  deleteService,
		})
	})
	return c.controller, c.controllersErr
}

func (c *MainLambdaContainer) GetMiddlewares() (trace, base, errorMw middlewares.Middleware) {
	c.middlewaresOnce.Do(func() {
		c.traceMiddleware = middlewares.TraceMiddleware(c.logger)
		c.baseMiddleware = middlewares.BaseMiddleware(c.logger)
		c.errorMiddleware = middlewares.ErrorMiddleware(c.logger)
	})
	return c.traceMiddleware, c.baseMiddleware, c.errorMiddleware
}
