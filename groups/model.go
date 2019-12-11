package groups

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/peteclark-io/footie/bookings"
	"github.com/peteclark-io/footie/players"
)

const tableName = "player_groups"
const tableKey = "id"

var errGroupNotFound = errors.New("Group not found")

type Group struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	Players    []*players.Player   `json:"players,omitempty"`
	PlayerIDs  []string            `json:"playerIDs,omitempty"`
	Bookings   []*bookings.Booking `json:"bookings,omitempty"`
	BookingIDs []string            `json:"bookingIDs,omitempty"`
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

func GetGroups(db *dynamodb.DynamoDB) ([]*Group, error) {
	res, err := db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Limit:     aws.Int64(64),
	})

	if err != nil {
		return nil, err
	}

	if len(res.Items) == 0 {
		return nil, errGroupNotFound
	}

	groups := make([]*Group, 0)
	for _, item := range res.Items {
		gr, err := parseItem(db, item)
		if err != nil {
			return nil, err
		}
		groups = append(groups, gr)
	}
	return groups, nil
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

	return parseItem(db, res.Item)
}

func parseItem(db *dynamodb.DynamoDB, item map[string]*dynamodb.AttributeValue) (*Group, error) {
	gr := &Group{}
	err := dynamodbattribute.UnmarshalMap(item, gr)

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

	for _, id := range gr.BookingIDs {
		b, err := bookings.GetBooking(db, id)
		if err == bookings.ErrBookingNotFound {
			continue
		}

		if err != nil {
			return nil, err
		}

		gr.Bookings = append(gr.Bookings, b)
	}
	gr.BookingIDs = nil
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
		"playerIDs":  g.toPlayerIDs(),
		"bookingIDs": g.toBookingIDs(),
	}
}

func (g *Group) toPlayerIDs() *dynamodb.AttributeValue {
	vals := make([]*dynamodb.AttributeValue, 0)
	for _, p := range g.Players {
		vals = append(vals, &dynamodb.AttributeValue{S: aws.String(p.ID)})
	}
	return &dynamodb.AttributeValue{L: vals}
}

func (g *Group) toBookingIDs() *dynamodb.AttributeValue {
	vals := make([]*dynamodb.AttributeValue, 0)
	for _, b := range g.Bookings {
		vals = append(vals, &dynamodb.AttributeValue{S: aws.String(b.ID)})
	}
	return &dynamodb.AttributeValue{L: vals}
}
