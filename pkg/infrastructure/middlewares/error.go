package middlewares

import (
	"context"
	"encoding/json"

	errorCustom "github.com/Yolto7/api-candidates/pkg/domain/error"
	"github.com/Yolto7/api-candidates/pkg/domain/logger"
	"github.com/aws/aws-lambda-go/events"
)

func ErrorMiddleware(log logger.Logger) Middleware {
	return func(next LambdaHandlerFunc) LambdaHandlerFunc {
		return func(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
			resp, err := next(ctx, event)
			if err == nil {
				return resp, nil
			}

			log.Error(map[string]any{
				"error": err,
			})
			
			appErr := errorCustom.FromError(err)
			payload := map[string]any{
				"success": false,
				"message": appErr.Message,
				"code":    appErr.ErrorCode,
			}

			if appErr.Payload != nil {
				payload["payload"] = appErr.Payload
			}

			body, _ := json.Marshal(payload)
			return &events.APIGatewayProxyResponse{
				StatusCode:      appErr.HttpCode,
				IsBase64Encoded: false,
				Body:            string(body),
				Headers: map[string]string{
					"Content-Type":                "application/json",
					"Access-Control-Allow-Origin": "*",
				},
			}, nil
		}
	}
}
