package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types" // DEBUG
	"github.com/magomzr/magomzr-api/models"
)

var (
	tableName = "posts"
)

// To avoid returning all posts data, this method just maps
// the relevant fields using the Card struct.
func GetAllPosts(ctx context.Context, dynamoClient *dynamodb.Client) ([]models.Card, error) {
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

func GetPostById(ctx context.Context, dynamoClient *dynamodb.Client, id string) (*models.Post, error) {
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

// Create (or update) posts.
func SavePost(ctx context.Context, dynamoClient *dynamodb.Client, post *models.Post, isNew bool) (bool, error) {
	currentDatetime := time.Now().Format(time.RFC3339)

	if isNew {
		post.GenerateId()
		post.CreateDate = currentDatetime
	} else {
		post.ModifiedDate = currentDatetime
	}

	post.TagsToLower()

	// DEBUG
	fmt.Printf("Post before marshaling - ID: %s, Title: %s\n", post.ID, post.Title)

	item, err := attributevalue.MarshalMap(post)

	if err != nil {
		return false, fmt.Errorf("error marshaling post: %w", err)
	}

	// DEBUG
	if _, exists := item["id"]; !exists {
		return false, fmt.Errorf("marshaled item is missing 'id' key. Available keys: %v", getKeys(item))
	}

	_, err = dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})

	if err != nil {
		return false, fmt.Errorf("error saving post: %w", err)
	}

	return true, nil
}

// DEBUG
func getKeys(item map[string]types.AttributeValue) []string {
	keys := make([]string, 0, len(item))
	for k := range item {
		keys = append(keys, k)
	}
	return keys
}
