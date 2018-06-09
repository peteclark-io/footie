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
	"github.com/peteclark-io/footie/utils"
)

// Handler does alllllll the logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters[bookings.TableKey]
	if !ok {
		return utils.HTTPResponse("Please provide an 'id' in the path", http.StatusBadRequest), nil
	}

	b := bookings.Booking{}
	err := json.Unmarshal([]byte(request.Body), &b)
	if err != nil {
		return utils.HTTPResponse("Body should be application/json", http.StatusBadRequest), nil
	}

	if b.ID != id {
		return utils.HTTPResponse("Booking 'id' in the body should equal the 'id' in the path", http.StatusBadRequest), nil
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err = db.PutItem(&dynamodb.PutItemInput{
		Item:      b.ToItem(),
		TableName: aws.String(bookings.TableName),
	})

	if err != nil {
		return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
	}

	return utils.HTTPResponse("Saved booking", http.StatusOK), nil
}

func main() {
	lambda.Start(Handler)
}
