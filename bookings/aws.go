package bookings

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/peteclark-io/footie/aws"
)

const bookingsResource = "bookings"

type awsHandler struct {
	repository *Repository
}

func NewAWSHandler() lambda.Handler {
	return &awsHandler{repository: &Repository{}}
}

func (h *awsHandler) Configure(resource, method string) (interface{}, bool) {
	if bookingsResource != resource {
		return nil, false
	}

	return lambda.FunctionForCRUDMethod(h, method), true
}

func (h *awsHandler) Create(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	b, status, err := unmarshalBooking(strings.NewReader(request.Body))
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	status, err = h.repository.Create(b)
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	js, _ := json.Marshal(b)
	return events.APIGatewayProxyResponse{Body: string(js), StatusCode: http.StatusOK, Headers: map[string]string{"Content-Type": "application/json"}}, nil
}

func (h *awsHandler) Delete(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, status, err := lambda.CheckID(tableKey, request)
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	status, err = h.repository.Delete(id)

	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	return lambda.Response("Booking deleted", http.StatusAccepted), nil
}

func (h *awsHandler) Read(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, status, err := lambda.CheckID(tableKey, request)
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	b, status, err := h.repository.Read(id)
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	js, _ := json.Marshal(b)
	return events.APIGatewayProxyResponse{Body: string(js), StatusCode: http.StatusOK, Headers: map[string]string{"Content-Type": "application/json"}}, nil
}

func (h *awsHandler) Write(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, status, err := lambda.CheckID(tableKey, request)
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	b, status, err := unmarshalBooking(strings.NewReader(request.Body))
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	status, err = h.repository.Write(id, b)

	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	return lambda.Response("Saved booking", http.StatusOK), nil
}
