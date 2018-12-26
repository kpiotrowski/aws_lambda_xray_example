package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/kpiotrowski/aws_lambda_xray_example/simple/item"
	"github.com/satori/go.uuid"
	"time"
)

func Handler(ctx context.Context) (events.APIGatewayProxyResponse, error) {

	ctx, seg := xray.BeginSubsegment(ctx, "Generating ID")
	time.Sleep(time.Second*3)
	u1 := uuid.NewV4()
	itemID := item.ItemID{ID: aws.String(u1.String())}

	err := xray.AddMetadata(ctx, "generated_id", u1.String())
	seg.Close(err)

	itemBytes, _ := json.Marshal(itemID)

	return events.APIGatewayProxyResponse{
		Body:       string(itemBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
