package middlewares

import (
	"context"
	"fmt"
	"strings"

	"github.com/Yolto7/api-candidates/pkg/domain/constants"
	"github.com/Yolto7/api-candidates/pkg/domain/logger"
	"github.com/Yolto7/api-candidates/pkg/domain/trace"
	"github.com/Yolto7/api-candidates/pkg/infrastructure/utils"
	"github.com/aws/aws-lambda-go/events"
)

func TraceMiddleware(log logger.Logger) Middleware {
	return func(next LambdaHandlerFunc) LambdaHandlerFunc {
		return func(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
			traceID := extractOrGenerateTraceID(event.Headers)

			log.Info(fmt.Sprintf("TraceID: %s", traceID))

			ctx = trace.SetTraceID(ctx, traceID)

			return next(ctx, event)
		}
	}
}

func extractOrGenerateTraceID(headers map[string]string) string {
	for key, val := range headers {
		if strings.EqualFold(key, constants.HEADER_TRACE_ID) && strings.TrimSpace(val) != "" {
			return val
		}
	}
	return utils.GenerateUUID()
}
