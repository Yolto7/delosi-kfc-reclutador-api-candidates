package response

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func Success(status int, message string, data interface{}) (*events.APIGatewayProxyResponse, error) {
	res := map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	}
	body, _ := json.Marshal(res)
	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
	}, nil
}