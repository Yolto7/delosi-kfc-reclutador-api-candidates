package dynamo

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	instance *dynamodb.Client
	once     sync.Once
)

func GetClient(ctx context.Context) (*dynamodb.Client, error) {
	var err error
	once.Do(func() {
		cfg, cfgErr := config.LoadDefaultConfig(ctx)
		if cfgErr != nil {
			err = fmt.Errorf("failed to load config for DynamoDB: %w", cfgErr)
			return
		}
		
		instance = dynamodb.NewFromConfig(cfg)
	})

	return instance, err
}