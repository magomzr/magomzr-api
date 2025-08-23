package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/magomzr/magomzr-api/models"
)

func GetTags(ctx context.Context, dynamoClient *dynamodb.Client) (models.Tags, error) {
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
		return nil, fmt.Errorf("error querying DynamoDB: %w", err)
	}

	// Counting
	tagCounter := make(map[string]int)
	var posts []models.Post

	err = attributevalue.UnmarshalListOfMaps(result.Items, &posts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling posts: %w", err)
	}

	for _, post := range posts {
		for _, tag := range post.Tags {
			if tag != "" {
				tagCounter[strings.ToLower(tag)]++
			}
		}
	}

	return tagCounter, nil
}

// To avoid returning all posts data, this method just maps
// the relevant fields using the Card struct.
func GetPostsByTag(ctx context.Context, dynamoClient *dynamodb.Client, tag string) ([]models.Card, error) {
	tableName := "posts"

	draftFilter := expression.Name("isDraft").Equal(expression.Value(false))
	tagFilter := expression.Name("tags").Contains(tag)
	combinedFilter := expression.And(draftFilter, tagFilter)

	expr, err := expression.NewBuilder().WithFilter(combinedFilter).Build()
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
