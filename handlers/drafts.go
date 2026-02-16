package handlers

import (
	"context"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/magomzr/magomzr-api/models"
)

func GetDrafts(ctx context.Context, dynamoClient *dynamodb.Client) ([]models.Post, error) {
	filt := expression.Name("isDraft").Equal(expression.Value(true))
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

	var posts []models.Post
	err = attributevalue.UnmarshalListOfMaps(result.Items, &posts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling posts: %w", err)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreateDate > posts[j].CreateDate
	})

	return posts, nil
}
