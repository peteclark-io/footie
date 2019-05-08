package matches

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

func (r *Repository) Write(id string, m *Match) (int, error) {
	if m.ID != id {
		return http.StatusBadRequest, errors.New("Match 'id' in the body should equal the 'id' in the path")
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      m.ToItem(),
		TableName: aws.String(tableName),
	})

	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	return http.StatusOK, nil
}

func (r *Repository) Read(id string) (*Match, int, error) {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	m, err := GetMatch(db, id)

	if err == ErrMatchNotFound {
		return nil, http.StatusNotFound, errors.New("Match not found")
	}

	if err != nil {
		return m, http.StatusServiceUnavailable, err
	}

	return m, 0, nil
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

func (r *Repository) Create(m *Match) (int, error) {
	if m.ID != "" {
		return http.StatusBadRequest, errors.New("New matches should not contain an 'id' field")
	}

	m.ID = ids.NewID()

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      m.ToItem(),
		TableName: aws.String(tableName),
	})

	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	return http.StatusOK, nil
}
