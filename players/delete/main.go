package main

import (
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

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err := db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(players.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			players.PlayerKey: {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
	}

	return utils.HTTPResponse("Player deleted", http.StatusAccepted), nil
}

func main() {
	lambda.Start(Handler)
}
