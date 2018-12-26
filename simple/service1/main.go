package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/kpiotrowski/aws_lambda_xray_example/simple/item"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
	"os"
)

//Handler is function executed by lambda engine
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	sess := session.Must(session.NewSession())
	snsClient := sns.New(sess)
	dynamoClient := dynamodb.New(sess)
	xray.AWS(snsClient.Client)
	xray.AWS(dynamoClient.Client)

	itemData := item.Item{}
	err := json.Unmarshal([]byte(request.Body), &itemData)
	if err != nil {
		return nil, err
	}

	itemID, err := get_id(ctx)
	if err != nil {
		return nil, err
	}
	itemData.Id = itemID.ID
	err = save_item(ctx, dynamoClient, itemData)
	if err != nil {
		return nil, err
	}

	sns_notify(ctx, snsClient, "Created new Item", fmt.Sprintf("Created new item with ID: %s", *itemID.ID))
	itemIdBytes, _ := json.Marshal(itemID)

	return &events.APIGatewayProxyResponse{
		Body:       string(itemIdBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}

func sns_notify(ctx context.Context, snsClient *sns.SNS, title, message string) {
	input := &sns.PublishInput{
		Subject: aws.String(title),
		Message: aws.String(message),
		TopicArn: aws.String(os.Getenv("NOTIFICATION_TOPIC")),
	}
	snsClient.PublishWithContext(ctx, input)
}

func save_item(ctx context.Context, dynamoClient *dynamodb.DynamoDB, itemData item.Item) error {
	data, err := dynamodbattribute.MarshalMap(itemData)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("ITEMS_TABLE")),
		Item: data,
	}
	_, err = dynamoClient.PutItemWithContext(ctx, input)
	if err != nil {
		err = xray.AddError(ctx, err)
	}
	return err
}

func get_id(ctx context.Context) (*item.ItemID, error) {
	httpClient := xray.Client(http.DefaultClient)
	resp, err := ctxhttp.Get(
		ctx,
		httpClient,
		fmt.Sprintf("https://%s/%s/id_generator",os.Getenv("API_HOST"), os.Getenv("API_STAGE")),
	)

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	itemID := item.ItemID{}
	err = json.Unmarshal(body, &itemID)

	return &itemID, err
}