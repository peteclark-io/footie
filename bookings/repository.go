package bookings

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/peteclark-io/footie/ids"
	"github.com/peteclark-io/footie/matches"
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

	if err == ErrBookingNotFound {
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

	if b.Cadence == "" {
		return http.StatusBadRequest, errors.New("Bookings should specify the cadence for matches, i.e. 168h for weekly")
	}

	if b.Start == nil {
		return http.StatusBadRequest, errors.New("New bookings require a start date")
	}

	if b.Length <= 0 {
		return http.StatusBadRequest, errors.New("New bookings require a non-zero length")
	}

	dur, err := time.ParseDuration(b.Cadence)
	if err != nil {
		return http.StatusBadRequest, errors.New("Cadence must be a valid duration")
	}

	b.ID = ids.NewID()
	b.cadenceDuration = dur
	b.Cadence = dur.String()

	matchRepo := &matches.Repository{}
	for i := 1; i <= b.Length; i++ {
		matchDay := b.Start.Add(b.cadenceDuration * time.Duration(i))
		m := &matches.Match{
			Name:     fmt.Sprintf("Game %v", i),
			StartsAt: matchDay,
			Group:    b.Group,
			Booking:  b.ID,
		}

		status, err := matchRepo.Create(m)
		if err != nil {
			return status, err
		}

		b.Matches = append(b.Matches, m)
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)
	_, err = db.PutItem(&dynamodb.PutItemInput{
		Item:      b.ToItem(),
		TableName: aws.String(tableName),
	})

	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	return http.StatusOK, nil
}
