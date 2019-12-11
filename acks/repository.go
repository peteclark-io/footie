package acks

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

func (r *Repository) Write(id string, a *Acknowledgement) (int, error) {
	if a.ID != id {
		return http.StatusBadRequest, errors.New("Acknowledgement 'id' in the body should match the 'id' in the path")
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      a.ToItem(),
		TableName: aws.String(tableName),
	})

	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	return http.StatusOK, nil
}

func (r *Repository) Read(id string) (*Acknowledgement, int, error) {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	b, err := GetAcknowledgement(db, id)

	if err == ErrAcknowledgementNotFound {
		return nil, http.StatusNotFound, errors.New("Acknowledgement not found")
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

func (r *Repository) Create(a *Acknowledgement) (int, error) {
	if a.ID != "" {
		return http.StatusBadRequest, errors.New("New acknowledgements should not contain an 'id' field")
	}

	a.ID = ids.NewID()

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)
	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      a.ToItem(),
		TableName: aws.String(tableName),
	})

	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	return http.StatusOK, nil
}
