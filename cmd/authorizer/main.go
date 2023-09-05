package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/kusipay/api-go-auth/middleware"
	"github.com/mefellows/vesper"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func errorResponse(err error) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: false,
		Context: map[string]interface{}{
			"message": err.Error(),
		},
	}, nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	authorizationValue, ok := event.Headers["authorization"]
	if !ok {
		return errorResponse(fmt.Errorf("authorization token not found"))
	}

	if !strings.HasPrefix(authorizationValue, "Bearer ") {
		return errorResponse(fmt.Errorf("invalid authorization token"))
	}

	tokenString := strings.TrimPrefix(authorizationValue, "Bearer ")

	region := os.Getenv("REGION")
	userPoolId := os.Getenv("USER_POOL_ID")

	jwksUri := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolId)

	cache := jwk.NewCache(ctx)

	err := cache.Register(jwksUri)
	if err != nil {
		return errorResponse(err)
	}

	set, err := cache.Get(ctx, jwksUri)
	if err != nil {
		return errorResponse(err)
	}

	token, err := jwt.Parse([]byte(tokenString), jwt.WithKeySet(set), jwt.WithValidate(true))
	if err != nil {
		return errorResponse(err)
	}

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context:      token.PrivateClaims(),
	}, nil
}

func main() {
	v := vesper.New(Handler).Use(middleware.LogMiddleware())

	v.Start()
}
