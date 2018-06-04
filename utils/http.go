package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type Msg struct {
	Msg string `json:"msg"`
}

func HTTPResponse(msg string, status int) events.APIGatewayProxyResponse {
	b, _ := json.Marshal(&Msg{Msg: msg})
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: status, Headers: map[string]string{"Content-Type": "application/json"}}
}
