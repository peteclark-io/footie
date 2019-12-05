package acks

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

const tableName = "acknowledgements"
const tableKey = "id"

var ErrAcknowledgementNotFound = errors.New(`Acknowledgement not found`)

type Acknowledgement struct {
	ID     string `json:"id"`
	Group  string `json:"group"`
	Player string `json:"player"`
	Match  string `json:"match"`
	Value  bool   `json:"value"`
}

func unmarshalAcknowledgement(body io.Reader) (*Acknowledgement, int, error) {
	a := Acknowledgement{}

	dec := json.NewDecoder(body)
	err := dec.Decode(&a)

	if err != nil {
		log.WithError(err).Error("Failed to unmarshal acknowledgement body")
		return nil, http.StatusBadRequest, errors.New("Body should be application/json")
	}

	return &a, 0, nil
}

func GetAcknowledgement(db *dynamodb.DynamoDB, id string) (*Acknowledgement, error) {
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
		return nil, ErrAcknowledgementNotFound
	}

	a := &Acknowledgement{}
	err = dynamodbattribute.UnmarshalMap(res.Item, a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Acknowledgement) ToItem() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(a.ID),
		},
		"group": {
			S: aws.String(a.Group),
		},
		"match": {
			S: aws.String(a.Match),
		},
		"player": {
			S: aws.String(a.Player),
		},
		"value": {
			BOOL: aws.Bool(a.Value),
		},
	}
}
