package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/peteclark-io/footie/ids"
	"github.com/peteclark-io/footie/players"
	"github.com/peteclark-io/footie/utils"
)

// Handler does alllllll the logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	pl := players.Player{}
	err := json.Unmarshal([]byte(request.Body), &pl)
	if err != nil {
		return utils.HTTPResponse("Body should be application/json", http.StatusBadRequest), nil
	}

	if pl.ID != "" {
		return utils.HTTPResponse("New players should not contain an 'id' field", http.StatusBadRequest), nil
	}

	pl.ID = ids.NewID()

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err = db.PutItem(&dynamodb.PutItemInput{
		Item:      pl.ToItem(),
		TableName: aws.String(players.TableName),
	})

	if err != nil {
		return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
	}

	b, _ := json.Marshal(pl)
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: http.StatusOK, Headers: map[string]string{"Content-Type": "application/json"}}, nil
}

func main() {
	lambda.Start(Handler)
}
