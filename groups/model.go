package groups

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/peteclark-io/footie/players"
)

const tableName = "player_groups"
const tableKey = "id"

var errGroupNotFound = errors.New("Group not found")

type Group struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Players   []*players.Player `json:"players"`
	PlayerIDs []string          `json:"playerIDs,omitempty"`
}

func unmarshalGroup(body io.Reader) (*Group, int, error) {
	gr := Group{}

	dec := json.NewDecoder(body)
	err := dec.Decode(&gr)

	if err != nil {
		return nil, http.StatusBadRequest, errors.New("Body should be application/json")
	}

	return &gr, 0, nil
}

func GetGroup(db *dynamodb.DynamoDB, id string) (*Group, error) {
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
		return nil, errGroupNotFound
	}

	gr := &Group{}
	err = dynamodbattribute.UnmarshalMap(res.Item, gr)

	if err != nil {
		return nil, err
	}

	for _, id := range gr.PlayerIDs {
		pl, err := players.GetPlayer(db, id)
		if err == players.ErrPlayerNotFound {
			continue
		}

		if err != nil {
			return nil, err
		}

		gr.Players = append(gr.Players, pl)
	}

	gr.PlayerIDs = nil
	return gr, nil
}

func (g *Group) ToItem() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(g.ID),
		},
		"name": {
			S: aws.String(g.Name),
		},
		"playerIDs": g.toPlayerIDs(),
	}
}

func (g *Group) toPlayerIDs() *dynamodb.AttributeValue {
	vals := make([]*dynamodb.AttributeValue, 0)
	for _, p := range g.Players {
		vals = append(vals, &dynamodb.AttributeValue{S: aws.String(p.ID)})
	}
	return &dynamodb.AttributeValue{L: vals}
}
