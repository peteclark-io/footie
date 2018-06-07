package players

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const TableName = "players"
const TableKey = "id"

var (
	ErrPlayerNotFound = errors.New("player not found")
)

type Player struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
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
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			TableKey: {
				S: aws.String(id),
			},
		},
	})

	if len(res.Item) == 0 {
		return nil, ErrPlayerNotFound
	}

	if err != nil {
		return nil, err
	}

	pl := &Player{}
	err = dynamodbattribute.UnmarshalMap(res.Item, pl)

	if err != nil {
		return nil, err
	}

	return pl, nil
}
