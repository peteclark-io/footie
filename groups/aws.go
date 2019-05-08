package groups

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/peteclark-io/footie/aws"
)

type awsHandler struct {
	repository *Repository
}

func NewAWSHandler() lambda.Handler {
	return &awsHandler{repository: &Repository{}}
}

func (h *awsHandler) Create(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	pl, status, err := unmarshalGroup(strings.NewReader(request.Body))
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	status, err = h.repository.Create(pl)
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	b, _ := json.Marshal(pl)
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: http.StatusOK, Headers: map[string]string{"Content-Type": "application/json"}}, nil
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

	return lambda.Response("Group deleted", http.StatusAccepted), nil
}

func (h *awsHandler) Read(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, status, err := lambda.CheckID(tableKey, request)
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	pl, status, err := h.repository.Read(id)
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	b, _ := json.Marshal(pl)
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: http.StatusOK, Headers: map[string]string{"Content-Type": "application/json"}}, nil
}

func (h *awsHandler) Write(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, status, err := lambda.CheckID(tableKey, request)
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	m, status, err := unmarshalGroup(strings.NewReader(request.Body))
	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	status, err = h.repository.Write(id, m)

	if err != nil {
		return lambda.Response(err.Error(), status), nil
	}

	return lambda.Response("Saved group", http.StatusOK), nil
}
