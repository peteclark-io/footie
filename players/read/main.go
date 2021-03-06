package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/peteclark-io/footie/players"
	"github.com/peteclark-io/footie/utils"
)

// Handler does alllllll the logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters[players.TableKey]
	if !ok {
		return events.APIGatewayProxyResponse{Body: "Please provide an 'id'", StatusCode: http.StatusBadRequest}, nil
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	pl, err := players.GetPlayer(db, id)

	if err == players.ErrPlayerNotFound {
		return utils.HTTPResponse("Player not found", http.StatusNotFound), nil
	}

	if err != nil {
		return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
	}

	b, _ := json.Marshal(pl)
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: http.StatusOK, Headers: map[string]string{"Content-Type": "application/json"}}, nil
}

func main() {
	lambda.Start(Handler)
}
