package acks

import (
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/peteclark-io/footie/aws"
	"github.com/sirupsen/logrus"
)

var addressRegex = regexp.MustCompile(`in\+([a-z0-9]{9})([a-z0-9]{9})([a-z0-9]{9})@shoreditch\.football`)

const ackResource = "acknowledgements"

const (
	IN  = "IN"
	OUT = "OUT"
)

type awsHandler struct {
	repository *Repository
}

func NewAWSHandler() lambda.Handler {
	return &awsHandler{repository: &Repository{}}
}

func (h *awsHandler) Configure(resource, method string) (interface{}, bool) {
	if ackResource != resource {
		return nil, false
	}

	return h.Receive, true
}

func (h *awsHandler) Receive(request events.SimpleEmailEvent) error {
	for _, msg := range request.Records {
		if len(msg.SES.Mail.CommonHeaders.To) > 1 {
			logrus.WithField("to", msg.SES.Mail.CommonHeaders.To).Warn("More than one recipient")
			continue
		}

		to := msg.SES.Mail.CommonHeaders.To[0]
		matches := addressRegex.FindStringSubmatch(to)
		if len(matches) != 4 {
			logrus.WithField("to", to).Warn("Recipient does not match expected labelling")
			continue
		}

		subject := msg.SES.Mail.CommonHeaders.Subject
		val := strings.Contains(subject, IN)
		if !val && !strings.Contains(subject, OUT) {
			logrus.WithField("subject", subject).Warn("Subject is neither IN nor OUT")
			continue
		}

		status, err := h.repository.Create(&Acknowledgement{Group: matches[1], Match: matches[2], Player: matches[3], Value: true})
		if err != nil {
			logrus.WithError(err).WithField("status", status).Error("Failed to write ack")
		}
	}
	return nil
}
