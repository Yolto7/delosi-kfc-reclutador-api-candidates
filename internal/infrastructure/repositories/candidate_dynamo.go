package repositories

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/Yolto7/api-candidates/internal/domain/entities"
	"github.com/Yolto7/api-candidates/internal/domain/repositories"
	errorCustom "github.com/Yolto7/api-candidates/pkg/domain/error"
	"github.com/Yolto7/api-candidates/pkg/domain/logger"
	"github.com/Yolto7/api-candidates/pkg/infrastructure/utils"
)

type CandidateDynamoRepository struct {
	logger logger.Logger
	client *dynamodb.Client
	table  string
}

func NewCandidateDynamoRepository(logger logger.Logger, client *dynamodb.Client, table string) *CandidateDynamoRepository {
	return &CandidateDynamoRepository{
		logger: logger,
		client: client,
		table:  table,
	}
}

func (r *CandidateDynamoRepository) GetByID(ctx context.Context, id string) (*entities.Candidate, error) {
	res, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		r.logger.Error(utils.NewSafeError(err, "Error in CandidateRepository.GetByID: Failed to get candidate"))
		return nil, errorCustom.NewError(errorCustom.BAD_REQUEST, "Failed to get candidate data", "DATABASE_ERROR")
	}
	if len(res.Item) == 0 {
		return nil, nil
	}

	var candidate entities.Candidate
	if err := attributevalue.UnmarshalMap(res.Item, &candidate); err != nil {
		r.logger.Error(utils.NewSafeError(err, "Error in CandidateRepository.GetByID: Failed to unmarshal candidate"))
		return nil, errorCustom.NewError(errorCustom.BAD_REQUEST, "Failed to unmarshal building", "DATABASE_ERROR")
	}

	return &candidate, nil
}

func (r *CandidateDynamoRepository) GetByCompositeKey(ctx context.Context, compositeKey string) (*entities.Candidate, error) {
	res, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
		IndexName:              aws.String("GSI-Candidates-CompositeKey"),
		KeyConditionExpression: aws.String("compositeKey = :ck"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ck": &types.AttributeValueMemberS{Value: compositeKey},
		},
	})
	if err != nil {
		r.logger.Error(utils.NewSafeError(err, "Error in CandidateRepository.GetByCompositeKey: Failed to query candidate"))
		return nil, errorCustom.NewError(errorCustom.BAD_REQUEST, "Failed to get candidate data", "DATABASE_ERROR")
	}

	if len(res.Items) == 0 {
		return nil, nil
	}

	var candidate entities.Candidate
	if err := attributevalue.UnmarshalMap(res.Items[0], &candidate); err != nil {
		r.logger.Error(utils.NewSafeError(err, "Error in CandidateRepository.GetByCompositeKey: Failed to unmarshal candidate"))
		return nil, errorCustom.NewError(errorCustom.BAD_REQUEST, "Failed to unmarshal candidate", "DATABASE_ERROR")
	}

	return &candidate, nil
}

func (r *CandidateDynamoRepository) Create(ctx context.Context, candidate *entities.Candidate) error {
	item, err := attributevalue.MarshalMap(candidate)
	if err != nil {
		r.logger.Error(utils.NewSafeError(err, "Error in CandidateRepository.Create: Failed to marshal candidate"))
		return errorCustom.NewError(errorCustom.BAD_REQUEST, "Failed to marshal candidate", "DATABASE_ERROR")
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      item,
	})
	if err != nil {
		r.logger.Error(utils.NewSafeError(err, "Error in CandidateRepository.Create: Failed to create candidate"))
		return errorCustom.NewError(errorCustom.BAD_REQUEST, "Failed to create candidate", "DATABASE_ERROR")
	}

	return nil
}

func (r *CandidateDynamoRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if id == "" {
		return errorCustom.NewError(errorCustom.BAD_REQUEST, "ID is required for update", "VALIDATION_ERROR")
	}

	if len(updates) == 0 {
		return nil
	}

	exprAttrNames := make(map[string]string)
	exprAttrValues := make(map[string]types.AttributeValue)

	var updateExpr strings.Builder
	updateExpr.WriteString("SET ")

	fieldCount := 0
	for key, value := range updates {
		av, err := attributevalue.Marshal(value)
		if err != nil {
			return errorCustom.NewError(errorCustom.BAD_REQUEST, fmt.Sprintf("Failed to marshal field %s", key), "DATABASE_ERROR")
		}

		if av, ok := av.(*types.AttributeValueMemberNULL); ok && av.Value {
			continue
		}

		if fieldCount > 0 {
			updateExpr.WriteString(", ")
		}

		placeholderName := "#f" + strconv.Itoa(fieldCount)
		placeholderValue := ":v" + strconv.Itoa(fieldCount)

		updateExpr.WriteString(placeholderName)
		updateExpr.WriteString(" = ")
		updateExpr.WriteString(placeholderValue)

		exprAttrNames[placeholderName] = key
		exprAttrValues[placeholderValue] = av
		fieldCount++
	}

	if fieldCount == 0 {
		return nil
	}

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(r.table),
		Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: id}},
		UpdateExpression:          aws.String(updateExpr.String()),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
	})
	if err != nil {
		r.logger.Error(utils.NewSafeError(err, "Error in CandidateRepository.Update: Update failed"))
		return errorCustom.NewError(errorCustom.BAD_REQUEST, "Failed to update visit", "DATABASE_ERROR")
	}

	return nil
}

func (r *CandidateDynamoRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		r.logger.Error(utils.NewSafeError(err, "Error in CandidateRepository.Delete: Delete failed"))
		return errorCustom.NewError(errorCustom.BAD_REQUEST, "Failed to delete candidate", "DATABASE_ERROR")
	}

	return nil
}

var _ repositories.CandidateRepository = (*CandidateDynamoRepository)(nil)
