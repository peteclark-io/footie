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
	"github.com/peteclark-io/footie/groups"
	"github.com/peteclark-io/footie/players"
	"github.com/peteclark-io/footie/utils"
)

// Handler does alllllll the logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters[groups.TableKey]
	if !ok {
		return events.APIGatewayProxyResponse{Body: "Please provide an 'id'", StatusCode: http.StatusBadRequest}, nil
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	res, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(groups.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			groups.TableKey: {
				S: aws.String(id),
			},
		},
	})

	if len(res.Item) == 0 {
		return utils.HTTPResponse("Group not found", http.StatusNotFound), nil
	}

	if err != nil {
		return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
	}

	gr := groups.Group{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &gr)

	if err != nil {
		return utils.HTTPResponse("Failed to unmarshal response", http.StatusServiceUnavailable), nil
	}

	for _, id := range gr.PlayerIDs {
		pl, err := players.GetPlayer(db, id)
		if err == players.ErrPlayerNotFound {
			continue
		}

		if err != nil {
			return utils.HTTPResponse(err.Error(), http.StatusServiceUnavailable), nil
		}

		gr.Players = append(gr.Players, pl)
	}

	gr.PlayerIDs = nil
	b, _ := json.Marshal(gr)
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: http.StatusOK, Headers: map[string]string{"Content-Type": "application/json"}}, nil
}

func main() {
	lambda.Start(Handler)
}
