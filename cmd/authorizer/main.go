package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/kusipay/api-go-auth/middleware"
	"github.com/mefellows/vesper"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func errorResponse(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       err.Error(),
	}, nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	tokenString := event.Headers["Authorization"]

	region := os.Getenv("REGION")
	userPoolId := os.Getenv("USER_POOL_ID")

	jwksUri := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolId)

	cache := jwk.NewCache(ctx)

	_ = cache.Register(jwksUri)

	set, err := cache.Get(ctx, jwksUri)
	if err != nil {
		return errorResponse(err)
	}

	token, err := jwt.Parse([]byte(tokenString), jwt.WithKeySet(set), jwt.WithValidate(true))
	if err != nil {
		return errorResponse(err)
	}

	body, _ := json.Marshal(token.PrivateClaims())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	v := vesper.New(Handler).Use(middleware.LogMiddleware())

	v.Start()
}
