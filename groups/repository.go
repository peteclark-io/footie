package groups

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

func (r *Repository) Write(id string, gr *Group) (int, error) {
	if gr.ID != id {
		return http.StatusBadRequest, errors.New("Group 'id' in the body should match the 'id' in the path")
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      gr.ToItem(),
		TableName: aws.String(tableName),
	})

	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	return http.StatusOK, nil
}

func (r *Repository) Read(id string) (*Group, int, error) {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	gr, err := GetGroup(db, id)

	if err == errGroupNotFound {
		return nil, http.StatusNotFound, errors.New("Group not found")
	}

	if err != nil {
		return gr, http.StatusServiceUnavailable, err
	}

	return gr, 0, nil
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

func (r *Repository) Create(gr *Group) (int, error) {
	if gr.ID != "" {
		return http.StatusBadRequest, errors.New("New groups should not contain an 'id' field")
	}

	gr.ID = ids.NewID()

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      gr.ToItem(),
		TableName: aws.String(tableName),
	})

	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	return http.StatusOK, nil
}
