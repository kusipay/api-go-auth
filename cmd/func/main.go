package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mefellows/vesper"
)

// Body
type Body struct {
	Error   bool   `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"Message"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	body := new(Body)
	_ = json.Unmarshal([]byte(event.Body), body)

	message, _ := json.Marshal(map[string]any{
		"message": body.Message,
	})

	if body.Error {
		return nil, errors.New(body.Message)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: body.Status,
		Body:       string(message),
	}, nil
}

func main() {
	v := vesper.New(Handler)

	v.Start()
}
