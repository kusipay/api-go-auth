package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mefellows/vesper"
)

type Jwk struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	Use string `json:"use"`
}

type JwkResponse struct {
	Keys []Jwk `json:"keys"`
}

func errorResponse(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       err.Error(),
	}, nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", "us-east-1", "us-east-1_VwhhbfJUF"))
	if err != nil {
		return errorResponse(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errorResponse(err)
	}

	jwks := new(JwkResponse)

	if err = json.Unmarshal(body, jwks); err != nil {
		return errorResponse(err)
	}

	tokenString := event.Headers["Authorization"]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		for _, key := range jwks.Keys {
			if token.Header["kid"] == key.Kid {
				return jwt.ParseRSAPublicKeyFromPEM([]byte(key.N))
			}
		}

		return nil, fmt.Errorf("no key found")
	})
	if err != nil {
		return errorResponse(err)
	}

	if !token.Valid {
		return errorResponse(fmt.Errorf("invalid token"))
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	v := vesper.New(Handler)

	v.Start()
}
