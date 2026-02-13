package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Yolto7/api-candidates/internal/infrastructure/container"
	"github.com/Yolto7/api-candidates/pkg/infrastructure/middlewares"
	"github.com/Yolto7/api-candidates/pkg/infrastructure/utils"
)

const PREFIX = "candidates"

var (
	initStart     time.Time
	initErr       error
	
	routes        map[string]map[string]middlewares.LambdaHandlerFunc
)

func init() {
	initStart = time.Now()
	ctx := context.Background()

	// Usar el nuevo container específico para main lambda
	var err error
	mainContainer, err := container.NewMainLambdaContainer(ctx)
	if err != nil {
		initErr = err
		return
	}

	// Configurar middlewares ya resueltos
	trace, base, errorMw := mainContainer.GetMiddlewares()
	baseMiddlewares := []middlewares.Middleware{trace, base, errorMw}

	controller, err := mainContainer.GetCandidateController()
	if err != nil {
		initErr = err
		return
	}

	// Configurar rutas con el controller específico	
	routes = map[string]map[string]middlewares.LambdaHandlerFunc{
		http.MethodGet: {
			PREFIX + "/{id}": middlewares.ChainMiddlewares(controller.GetByID, baseMiddlewares...),
		},
		http.MethodPost: {
			PREFIX + "/": middlewares.ChainMiddlewares(controller.Create, baseMiddlewares...),
		},
		http.MethodDelete: {
			PREFIX + "/{id}": middlewares.ChainMiddlewares(controller.Delete, baseMiddlewares...),
		},
	}

	mainContainer.Logger().Info(fmt.Sprintf("Main lambda init completed in %v", time.Since(initStart)))
}

func main() {
	lambda.Start(func(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		if initErr != nil {
			return nil, initErr
		}

		utilsRoutes := make(map[string]map[string]utils.LambdaHandlerFunc)
		for method, methodRoutes := range routes {
			utilsRoutes[method] = make(map[string]utils.LambdaHandlerFunc)
			for pattern, handler := range methodRoutes {
				utilsRoutes[method][pattern] = utils.LambdaHandlerFunc(handler)
			}
		}

		return utils.HandleRoutes(ctx, event, utilsRoutes)
	})
}
