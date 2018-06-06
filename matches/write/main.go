package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/peteclark-io/footie/matches"
	"github.com/peteclark-io/footie/utils"
)

// Handler does alllllll the logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters[matches.TableKey]
	if !ok {
		return utils.HTTPResponse("Please provide an 'id' in the path", http.StatusBadRequest), nil
	}

	m := matches.Match{}
	err := json.Unmarshal([]byte(request.Body), &m)
	if err != nil {
		return utils.HTTPResponse("Body should be application/json", http.StatusBadRequest), nil
	}

	if m.ID != id {
		return utils.HTTPResponse("Match 'id' in the body should equal the 'id' in the path", http.StatusBadRequest), nil
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err = db.PutItem(&dynamodb.PutItemInput{
		Item:      m.ToItem(),
		TableName: aws.String(matches.TableName),
	})

	if err != nil {
		return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
	}

	return utils.HTTPResponse("Saved match", http.StatusOK), nil
}

func main() {
	lambda.Start(Handler)
}
