package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

//Handler is function executed by lambda engine
func Handler(ctx context.Context) (string, error) {
	return "", nil
}

func main() {
	lambda.Start(Handler)
}
