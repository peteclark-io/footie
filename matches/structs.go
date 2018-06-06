package matches

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const TableName = "matches"
const TableKey = "id"

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
