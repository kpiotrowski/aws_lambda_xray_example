package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/kpiotrowski/aws_lambda_xray_example/simple/item"
	"os"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	itemID := request.PathParameters["id"]
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(itemID),
			},
		},
		TableName: aws.String(os.Getenv("ITEMS_TABLE")),
	}
	dynamoClient := dynamodb.New(session.Must(session.NewSession()))
	xray.AWS(dynamoClient.Client)

	resp, err := dynamoClient.GetItemWithContext(ctx, input)
	if err != nil {
		_ = xray.AddError(ctx, err)
		return nil, err
	}

	itemData := item.Item{}
	err = dynamodbattribute.UnmarshalMap(resp.Item, &itemData)
	if err != nil {
		return nil, err
	} else if itemData.Id == nil {
		return &events.APIGatewayProxyResponse{StatusCode: 404}, nil
	}

	itemBytes, _ := json.Marshal(itemData)

	return &events.APIGatewayProxyResponse{
		Body:       string(itemBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
