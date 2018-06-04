package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/peteclark-io/footie/players"
	"github.com/peteclark-io/footie/utils"
)

// Handler does alllllll the logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters[players.PlayerKey]
	if !ok {
		return utils.HTTPResponse("Please provide an 'id'", http.StatusBadRequest), nil
	}

	pl := players.Player{}
	err := json.Unmarshal([]byte(request.Body), &pl)
	if err != nil {
		return utils.HTTPResponse("Body should be application/json", http.StatusBadRequest), nil
	}

	if pl.ID != id {
		return utils.HTTPResponse("Player 'id' in the body should match the 'id' in the path", http.StatusBadRequest), nil
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err = db.PutItem(&dynamodb.PutItemInput{
		Item:      pl.ToItem(),
		TableName: aws.String(players.TableName),
	})

	if err != nil {
		return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
	}

	return utils.HTTPResponse("Saved player", http.StatusOK), nil
}

func main() {
	lambda.Start(Handler)
}
