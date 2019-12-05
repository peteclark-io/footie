package matches

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/peteclark-io/footie/players"
)

const tableName = "matches"
const tableKey = "id"

var ErrMatchNotFound = errors.New("match not found")

type Match struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	StartsAt  time.Time         `json:"startsAt"`
	Group     string            `json:"group"`
	Booking   string            `json:"booking"`
	PlayerIDs []string          `json:"playerIDs,omitempty"`
	Players   []*players.Player `json:"players,omitempty"`
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
		"name": {
			S: aws.String(m.Name),
		},
		"startsAt": {
			S: aws.String(m.StartsAt.UTC().Format(time.RFC3339Nano)),
		},
		"group": {
			S: aws.String(m.Group),
		},
		"booking": {
			S: aws.String(m.Booking),
		},
		"playerIDs": m.toPlayerIDs(),
	}
}

func (m *Match) toPlayerIDs() *dynamodb.AttributeValue {
	vals := make([]*dynamodb.AttributeValue, 0)
	for _, p := range m.Players {
		vals = append(vals, &dynamodb.AttributeValue{S: aws.String(p.ID)})
	}
	return &dynamodb.AttributeValue{L: vals}
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

	for _, id := range m.PlayerIDs {
		p, err := players.GetPlayer(db, id)
		if err == players.ErrPlayerNotFound {
			continue
		}

		if err != nil {
			return nil, err
		}
		m.Players = append(m.Players, p)
	}

	m.PlayerIDs = nil
	return m, nil
}
