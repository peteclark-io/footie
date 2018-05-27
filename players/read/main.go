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
)

const tableName = "players"
const playerKey = "id"

type Player struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// Handler does alllllll the logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters[playerKey]
	if !ok {
		return events.APIGatewayProxyResponse{Body: "Please provide an 'id'", StatusCode: http.StatusBadRequest}, nil
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	res, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			playerKey: {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusServiceUnavailable}, nil
	}

	pl := Player{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &pl)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusServiceUnavailable}, nil
	}

	b, _ := json.Marshal(pl)
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(Handler)
}
