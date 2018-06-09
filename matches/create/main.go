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
	"github.com/peteclark-io/footie/matches"
	"github.com/peteclark-io/footie/utils"
)

// Handler does alllllll the logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	m := matches.Match{}
	err := json.Unmarshal([]byte(request.Body), &m)
	if err != nil {
		return utils.HTTPResponse("Body should be application/json", http.StatusBadRequest), nil
	}

	if m.ID != "" {
		return utils.HTTPResponse("New matches should not contain an 'id' field", http.StatusBadRequest), nil
	}

	m.ID = ids.NewID()

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err = db.PutItem(&dynamodb.PutItemInput{
		Item:      m.ToItem(),
		TableName: aws.String(matches.TableName),
	})

	if err != nil {
		return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
	}

	b, _ := json.Marshal(m)
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: http.StatusOK, Headers: map[string]string{"Content-Type": "application/json"}}, nil
}

func main() {
	lambda.Start(Handler)
}
