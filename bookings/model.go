package bookings

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/peteclark-io/footie/matches"
)

const tableName = "group_bookings"
const tableKey = "id"

var errBookingNotFound = errors.New(`Booking not found`)

type Booking struct {
	ID       string           `json:"id"`
	Starts   *time.Time       `json:"starts"`
	Matches  []*matches.Match `json:"matches"`
	MatchIDs []string         `json:"matchIDs,omitempty"`
}

func unmarshalBooking(body io.Reader) (*Booking, int, error) {
	b := Booking{}

	dec := json.NewDecoder(body)
	err := dec.Decode(&b)

	if err != nil {
		return nil, http.StatusBadRequest, errors.New("Body should be application/json")
	}

	return &b, 0, nil
}

func GetBooking(db *dynamodb.DynamoDB, id string) (*Booking, error) {
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
		return nil, errBookingNotFound
	}

	b := &Booking{}
	err = dynamodbattribute.UnmarshalMap(res.Item, b)

	if err != nil {
		return nil, err
	}

	for _, id := range b.MatchIDs {
		m, err := matches.GetMatch(db, id)
		if err == matches.ErrMatchNotFound {
			continue
		}

		if err != nil {
			return nil, err
		}

		b.Matches = append(b.Matches, m)
	}

	b.MatchIDs = nil
	return b, nil
}

func (b *Booking) ToItem() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(b.ID),
		},
		"starts": {
			S: aws.String(b.Starts.UTC().Format(time.RFC3339Nano)),
		},
		"matchIDs": b.ToMatchIDs(),
	}
}

func (b *Booking) ToMatchIDs() *dynamodb.AttributeValue {
	vals := make([]*dynamodb.AttributeValue, 0)
	for _, m := range b.Matches {
		vals = append(vals, &dynamodb.AttributeValue{S: aws.String(m.ID)})
	}
	return &dynamodb.AttributeValue{L: vals}
}