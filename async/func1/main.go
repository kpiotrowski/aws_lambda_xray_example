package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

//Handler is function executed by lambda engine
func Handler(ctx context.Context) (string, error) {
	topicArn := os.Getenv("NotifyTopicArn")


	return topicArn, nil
}

func main() {
	lambda.Start(Handler)
}
