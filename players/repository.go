package players

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

func (r *Repository) Write(id string, pl *Player) (int, error) {
	if pl.ID != id {
		return http.StatusBadRequest, errors.New("Player 'id' in the body should match the 'id' in the path")
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      pl.ToItem(),
		TableName: aws.String(tableName),
	})

	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	return http.StatusOK, nil
}

func (r *Repository) Read(id string) (*Player, int, error) {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	pl, err := GetPlayer(db, id)

	if err == ErrPlayerNotFound {
		return nil, http.StatusNotFound, errors.New("Player not found")
	}

	if err != nil {
		return pl, http.StatusServiceUnavailable, err
	}

	return pl, 0, nil
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

func (r *Repository) Create(pl *Player) (int, error) {
	if pl.ID != "" {
		return http.StatusBadRequest, errors.New("New players should not contain an 'id' field")
	}

	pl.ID = ids.NewID()

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      pl.ToItem(),
		TableName: aws.String(tableName),
	})

	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	return http.StatusOK, nil
}
