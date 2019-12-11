package aws

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/peteclark-io/footie/resources"
)

type Handler interface {
	Configure(resource, method string) (interface{}, bool)
}

type CRUDHandler interface {
	Read(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	Write(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	Create(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	Delete(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func Response(msg string, status int) events.APIGatewayProxyResponse {
	b, _ := json.Marshal(&resources.Msg{Msg: msg})
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: status, Headers: map[string]string{"Content-Type": "application/json"}}
}

func CheckID(tableKey string, request events.APIGatewayProxyRequest) (string, int, error) {
	id, ok := request.PathParameters[tableKey]
	if !ok {
		return id, http.StatusBadRequest, errors.New("Please provide an 'id'")
	}
	return id, 0, nil
}

func FunctionForCRUDMethod(handler CRUDHandler, method string) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch method {
	case "create":
		return handler.Create
	case "write":
		return handler.Write
	case "read":
		return handler.Read
	case "delete":
		return handler.Delete
	}
	return nil
}
