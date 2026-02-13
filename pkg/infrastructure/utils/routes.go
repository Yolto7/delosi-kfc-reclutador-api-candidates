package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type LambdaHandlerFunc func(context.Context, events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

func MatchDynamicRoute(pattern, path string) (map[string]string, bool) {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	params := make(map[string]string)
	for i := range patternParts {
		pp := patternParts[i]
		cp := pathParts[i]

		if strings.HasPrefix(pp, "{") && strings.HasSuffix(pp, "}") {
			paramName := strings.Trim(pp, "{}")
			params[paramName] = cp
		} else if strings.ToLower(pp) != strings.ToLower(cp) {
			return nil, false
		}
	}
	return params, true
}

func GetSortedRoutePatterns(routeMap map[string]interface{}) []string {
	patterns := make([]string, 0, len(routeMap))
	for pattern := range routeMap {
		patterns = append(patterns, pattern)
	}

	sort.Slice(patterns, func(i, j int) bool {
		patternA := patterns[i]
		patternB := patterns[j]

		// 1. Rutas sin parámetros dinámicos primero
		hasParamsA := strings.Contains(patternA, "{")
		hasParamsB := strings.Contains(patternB, "{")

		if hasParamsA != hasParamsB {
			return !hasParamsA // Sin parámetros primero
		}

		// 2. Si ambas tienen o no tienen parámetros, más larga primero
		return len(patternA) > len(patternB)
	})

	return patterns
}

func HandleRoutes(ctx context.Context, event events.APIGatewayProxyRequest, routes map[string]map[string]LambdaHandlerFunc) (*events.APIGatewayProxyResponse, error) {
	method := event.HTTPMethod
	path := event.Path

	if methodRoutes, ok := routes[method]; ok {
		// Convertir map a interface{} para usar la función genérica
		routeMapInterface := make(map[string]interface{})
		for k := range methodRoutes {
			routeMapInterface[k] = nil
		}

		// Obtener patrones ordenados por especificidad
		sortedPatterns := GetSortedRoutePatterns(routeMapInterface)

		// Iterar en orden (las más específicas primero)
		for _, pattern := range sortedPatterns {
			handler := methodRoutes[pattern]

			// Intentar match exacto
			if pattern == path {
				return handler(ctx, event)
			}

			// Intentar match dinámico
			if params, matched := MatchDynamicRoute(pattern, path); matched {
				if event.PathParameters == nil {
					event.PathParameters = make(map[string]string)
				}
				for k, v := range params {
					event.PathParameters[k] = v
				}
				return handler(ctx, event)
			}
		}
	}

	// Route not found
	body, _ := json.Marshal(map[string]interface{}{
		"success": false,
		"message": "Route not found",
		"code":    "ROUTE_NOT_FOUND",
	})
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       string(body),
	}, nil
}