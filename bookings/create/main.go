package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/peteclark-io/footie/bookings"
	"github.com/peteclark-io/footie/ids"
	"github.com/peteclark-io/footie/utils"
)

// Handler does alllllll the logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	b := bookings.Booking{}
	err := json.Unmarshal([]byte(request.Body), &b)
	if err != nil {
		return utils.HTTPResponse(err.Error(), http.StatusBadRequest), nil
	}

	if b.ID != "" {
		return utils.HTTPResponse("New bookings should not contain an 'id' field", http.StatusBadRequest), nil
	}

	b.ID = ids.NewID()

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err = db.PutItem(&dynamodb.PutItemInput{
		Item:      b.ToItem(),
		TableName: aws.String(bookings.TableName),
	})

	if err != nil {
		return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
	}

	by, _ := json.Marshal(b)
	return events.APIGatewayProxyResponse{Body: string(by), StatusCode: http.StatusOK, Headers: map[string]string{"Content-Type": "application/json"}}, nil
}

func main() {
	lambda.Start(Handler)
}
