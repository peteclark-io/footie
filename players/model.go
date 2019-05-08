package players

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const tableName = "players"
const tableKey = "id"

var (
	// ErrPlayerNotFound occurs when we can't find the player you're after
	ErrPlayerNotFound = errors.New("player not found")
)

// Player contains player details
type Player struct {
	ID          string `json:"id"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

func unmarshalPlayer(body io.Reader) (*Player, int, error) {
	pl := Player{}

	dec := json.NewDecoder(body)
	err := dec.Decode(&pl)

	if err != nil {
		return nil, http.StatusBadRequest, errors.New("Body should be application/json")
	}

	return &pl, 0, nil
}

func (p *Player) ToItem() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(p.ID),
		},
		"email": {
			S: aws.String(p.Email),
		},
		"displayName": {
			S: aws.String(p.DisplayName),
		},
	}
}

func GetPlayer(db *dynamodb.DynamoDB, id string) (*Player, error) {
	res, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			tableKey: {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, ErrPlayerNotFound
	}

	pl := &Player{}
	err = dynamodbattribute.UnmarshalMap(res.Item, pl)

	if err != nil {
		return nil, err
	}

	return pl, nil
}
