package players

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const TableName = "players"
const TableKey = "id"

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
