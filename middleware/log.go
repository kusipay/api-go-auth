package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kusipay/api-go-auth/util"
	"github.com/mefellows/vesper"
)

// LogMiddleware middleware to log environments and events.
func LogMiddleware() vesper.Middleware {
	return func(next vesper.LambdaFunc) vesper.LambdaFunc {
		return func(ctx context.Context, event any) (any, error) {
			logEnvironments()
			logAny("Event", event)

			response, err := next(ctx, event)

			logAny("Response", response)
			logError(err)

			return response, err
		}
	}
}

func logEnvironments() {
	environments := os.Environ()

	envs := strings.Join(environments, "\r")

	util.Log("Environment |", envs)
}

func logAny(tag string, event any) {
	bytes, err := json.MarshalIndent(event, "", "  ")

	var text string
	if err != nil {
		text = fmt.Sprintf("%+v", event)
	} else {
		text = string(bytes)
	}

	util.Log(tag+" |", text)
}

func logError(err error) {
	if err != nil {
		util.Log("Error |", err.Error())
	}
}
