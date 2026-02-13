package middlewares

import (
	"context"

	"github.com/Yolto7/api-candidates/pkg/domain/logger"
	"github.com/aws/aws-lambda-go/events"
)

type LambdaHandlerFunc func(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

type Middleware func(LambdaHandlerFunc) LambdaHandlerFunc

func ChainMiddlewares(handler LambdaHandlerFunc, middlewares ...Middleware) LambdaHandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func BaseMiddleware(log logger.Logger) func(LambdaHandlerFunc) LambdaHandlerFunc {
	return func(next LambdaHandlerFunc) LambdaHandlerFunc {
		return func(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
			// Log completo del request
			log.Info(map[string]any{
				"msg":      "Incoming request",
				"path":     event.Path,
				"method":   event.HTTPMethod,
				"headers":  event.Headers,
				"query":    event.QueryStringParameters,
				"pathVars": event.PathParameters,
				"body":     event.Body,
			})

			resp, err := next(ctx, event)
			if err != nil {
				return resp, err
			}

			// Agregar headers de seguridad y CORS
			if resp.Headers == nil {
				resp.Headers = make(map[string]string)
			}

			// CORS
			resp.Headers["Access-Control-Allow-Origin"] = "*"
			resp.Headers["Content-Type"] = "application/json"

			// Seguridad
			resp.Headers["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains; preload"
			resp.Headers["Content-Security-Policy"] = "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self'; frame-ancestors 'self'"
			resp.Headers["X-Content-Type-Options"] = "nosniff"
			resp.Headers["X-Frame-Options"] = "DENY"
			resp.Headers["X-XSS-Protection"] = "1; mode=block"
			resp.Headers["Referrer-Policy"] = "no-referrer"

			return resp, nil
		}
	}
}