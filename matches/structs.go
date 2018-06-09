package matches

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const TableName = "matches"
const TableKey = "id"

var ErrMatchNotFound = errors.New("match not found")

type Match struct {
	ID string `json:"id"`
}

func (m *Match) ToItem() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(m.ID),
		},
	}
}

func GetMatch(db *dynamodb.DynamoDB, id string) (*Match, error) {
	res, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			TableKey: {
				S: aws.String(id),
			},
		},
	})

	if len(res.Item) == 0 {
		return nil, ErrMatchNotFound
	}

	if err != nil {
		return nil, err
	}

	m := &Match{}
	err = dynamodbattribute.UnmarshalMap(res.Item, m)

	if err != nil {
		return nil, err
	}

	return m, nil
}
