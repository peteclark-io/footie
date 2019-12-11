package tasks

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	lambda "github.com/peteclark-io/footie/aws"
	"github.com/peteclark-io/footie/emails"
	"github.com/peteclark-io/footie/groups"
	"github.com/peteclark-io/footie/matches"
)

const (
	Sender   = `"Monday Night Football" <in@shoreditch.football>`
	Subject  = "Are You In?"
	TextBody = `Are you in? If you are, reply "IN" to this email or "OUT" if you're not!`
)

const tasksResource = "tasks"

type awsHandler struct {
}

func NewAWSHandler() lambda.Handler {
	return &awsHandler{}
}

func (h *awsHandler) Configure(resource, method string) (interface{}, bool) {
	if tasksResource != resource {
		return nil, false
	}

	return h.EmailAcknowledgements, true
}

func (h *awsHandler) EmailAcknowledgements() error {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	gs, err := groups.GetGroups(db)
	if err != nil {
		return err
	}

	ms := make([]*matches.Match, 0)
	for _, gr := range gs {
		for _, b := range gr.Bookings {
			if b.End == nil {
				continue
			}

			if b.End.Before(time.Now()) {
				continue
			}

			m := b.NextMatchAfter(time.Now())
			if m == nil {
				continue
			}
			ms = append(ms, m)
		}
	}

	for _, m := range ms {
		processMatch(m)
	}

	return nil
}

func processMatch(match *matches.Match) error {
	for _, pl := range match.Players {
		ackEmail, err := emails.GenerateAcknowledgementEmail(pl, match)
		if err != nil {
			return err
		}

		producer := emails.NewEmailProducer()
		err = producer.ProduceEmail(Sender, pl.Email, Subject, ackEmail, TextBody)
		if err != nil {
			return err
		}
	}

	return nil
}
