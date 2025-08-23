package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/magomzr/magomzr-api/models"
)

func GetPostsByTag(ctx context.Context, dynamoClient *dynamodb.Client, tag string) ([]models.Post, error) {
	tableName := "posts"

	draftFilter := expression.Name("isDraft").Equal(expression.Value(false))
	tagFilter := expression.Name("tags").Contains(tag)
	combinedFilter := expression.And(draftFilter, tagFilter)

	expr, err := expression.NewBuilder().WithFilter(combinedFilter).Build()
	if err != nil {
		return nil, fmt.Errorf("error building expression: %w", err)
	}

	result, err := dynamoClient.Scan(ctx, &dynamodb.ScanInput{
		TableName:        &tableName,
		FilterExpression: expr.Filter(),
	})
	if err != nil {
		return nil, fmt.Errorf("error scanning DynamoDB: %w", err)
	}

	var posts []models.Post
	err = attributevalue.UnmarshalListOfMaps(result.Items, &posts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling posts: %w", err)
	}

	return posts, nil
}
