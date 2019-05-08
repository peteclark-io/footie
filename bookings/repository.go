package bookings

import (
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/peteclark-io/footie/ids"
)

type Repository struct {
}

func (r *Repository) Write(id string, b *Booking) (int, error) {
	if b.ID != id {
		return http.StatusBadRequest, errors.New("Booking 'id' in the body should match the 'id' in the path")
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      b.ToItem(),
		TableName: aws.String(tableName),
	})

	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	return http.StatusOK, nil
}

func (r *Repository) Read(id string) (*Booking, int, error) {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	b, err := GetBooking(db, id)

	if err == errBookingNotFound {
		return nil, http.StatusNotFound, errors.New("Booking not found")
	}

	if err != nil {
		return b, http.StatusServiceUnavailable, err
	}

	return b, 0, nil
}

func (r *Repository) Delete(id string) (int, error) {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err := db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			tableKey: {
				S: aws.String(id),
			},
		},
	})

	return http.StatusServiceUnavailable, err
}

func (r *Repository) Create(b *Booking) (int, error) {
	if b.ID != "" {
		return http.StatusBadRequest, errors.New("New bookings should not contain an 'id' field")
	}

	b.ID = ids.NewID()

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      b.ToItem(),
		TableName: aws.String(tableName),
	})

	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	return http.StatusOK, nil
}
