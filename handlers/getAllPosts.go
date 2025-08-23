package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/magomzr/magomzr-api/models"
)

// To avoid returning all posts data, this method just maps
// the relevant fields using the Card struct.
func GetAllPosts(ctx context.Context, dynamoClient *dynamodb.Client) ([]models.Card, error) {
	tableName := "posts"

	filt := expression.Name("isDraft").Equal(expression.Value(false))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, fmt.Errorf("error building expression: %w", err)
	}

	result, err := dynamoClient.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 &tableName,
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("error scanning DynamoDB: %w", err)
	}

	var cards []models.Card
	err = attributevalue.UnmarshalListOfMaps(result.Items, &cards)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling posts: %w", err)
	}

	return cards, nil
}
