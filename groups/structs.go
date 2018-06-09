package groups

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/peteclark-io/footie/players"
)

const TableName = "player_groups"
const TableKey = "id"

type Group struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Players   []*players.Player `json:"players"`
	PlayerIDs []string          `json:"playerIDs,omitempty"`
}

func (p *Group) ToItem() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(p.ID),
		},
		"name": {
			S: aws.String(p.Name),
		},
		"playerIDs": p.ToPlayerIDs(),
	}
}

func (g *Group) ToPlayerIDs() *dynamodb.AttributeValue {
	vals := make([]*dynamodb.AttributeValue, 0)
	for _, p := range g.Players {
		vals = append(vals, &dynamodb.AttributeValue{S: aws.String(p.ID)})
	}
	return &dynamodb.AttributeValue{L: vals}
}
