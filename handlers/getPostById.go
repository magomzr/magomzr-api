package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/magomzr/magomzr-api/models"
)

func GetPostById(ctx context.Context, dynamoClient *dynamodb.Client, id string) (*models.Post, error) {
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

	var posts []models.Post
	err = attributevalue.UnmarshalListOfMaps(result.Items, &posts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling post: %w", err)
	}

	// Posts
	var currentPostIndex int
	for i, post := range posts {
		if post.ID == id {
			currentPostIndex = i
			break
		}
		currentPostIndex = -1
	}

	if currentPostIndex == -1 {
		return nil, fmt.Errorf("post with ID %s not found", id)
	}

	// Get post by index

	currentPost := posts[currentPostIndex]

	if currentPostIndex > 0 {
		previousPost := posts[currentPostIndex-1]
		currentPost.Previous.ID = previousPost.ID
		currentPost.Previous.Title = previousPost.Title
	}

	if currentPostIndex < len(posts)-1 {
		nextPost := posts[currentPostIndex+1]
		currentPost.Next.ID = nextPost.ID
		currentPost.Next.Title = nextPost.Title
	}

	return &currentPost, nil
}
