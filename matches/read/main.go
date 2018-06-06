package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/peteclark-io/footie/matches"
	"github.com/peteclark-io/footie/utils"
)

// Handler does alllllll the logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters[matches.TableKey]
	if !ok {
		return utils.HTTPResponse("Please provide an 'id'", http.StatusBadRequest), nil
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	res, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(matches.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			matches.TableKey: {
				S: aws.String(id),
			},
		},
	})

	if len(res.Item) == 0 {
		return utils.HTTPResponse("Match not found", http.StatusNotFound), nil
	}

	if err != nil {
		return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
	}

	m := matches.Match{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &m)

	if err != nil {
		return utils.HTTPResponse("Failed to unmarshal response", http.StatusServiceUnavailable), nil
	}

	b, _ := json.Marshal(m)
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: http.StatusOK, Headers: map[string]string{"Content-Type": "application/json"}}, nil
}

func main() {
	lambda.Start(Handler)
}
