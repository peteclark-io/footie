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
	"github.com/peteclark-io/footie/bookings"
	"github.com/peteclark-io/footie/matches"
	"github.com/peteclark-io/footie/utils"
)

// Handler does alllllll the logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters[bookings.TableKey]
	if !ok {
		return utils.HTTPResponse("Please provide an 'id'", http.StatusBadRequest), nil
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	res, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(bookings.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			bookings.TableKey: {
				S: aws.String(id),
			},
		},
	})

	if len(res.Item) == 0 {
		return utils.HTTPResponse("Booking not found", http.StatusNotFound), nil
	}

	if err != nil {
		return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
	}

	b := bookings.Booking{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &b)

	if err != nil {
		return utils.HTTPResponse("Failed to unmarshal response", http.StatusServiceUnavailable), nil
	}

	for _, id := range b.MatchIDs {
		m, err := matches.GetMatch(db, id)
		if err == matches.ErrMatchNotFound {
			continue
		}

		if err != nil {
			return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
		}

		b.Matches = append(b.Matches, m)
	}

	b.MatchIDs = nil

	by, _ := json.Marshal(b)
	return events.APIGatewayProxyResponse{Body: string(by), StatusCode: http.StatusOK, Headers: map[string]string{"Content-Type": "application/json"}}, nil
}

func main() {
	lambda.Start(Handler)
}
