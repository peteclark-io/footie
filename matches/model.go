package matches

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const tableName = "matches"
const tableKey = "id"

var ErrMatchNotFound = errors.New("match not found")

type Match struct {
	ID string `json:"id"`
}

func unmarshalMatch(body io.Reader) (*Match, int, error) {
	m := Match{}

	dec := json.NewDecoder(body)
	err := dec.Decode(&m)

	if err != nil {
		return nil, http.StatusBadRequest, errors.New("Body should be application/json")
	}

	return &m, 0, nil
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
		return nil, ErrMatchNotFound
	}

	m := &Match{}
	err = dynamodbattribute.UnmarshalMap(res.Item, m)

	if err != nil {
		return nil, err
	}

	return m, nil
}
