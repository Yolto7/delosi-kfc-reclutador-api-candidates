package trace

import "context"

type contextKey string

const traceIDKey contextKey = "traceId"

func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		return "unknown-trace-id"
	}
	return traceID
}
