package bookings

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/peteclark-io/footie/matches"
)

const TableName = "group_bookings"
const TableKey = "id"

type Booking struct {
	ID       string           `json:"id"`
	Starts   *time.Time       `json:"starts"`
	Matches  []*matches.Match `json:"matches"`
	MatchIDs []string         `json:"matchIDs,omitempty"`
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
