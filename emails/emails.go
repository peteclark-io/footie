package emails

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/peteclark-io/footie/matches"
	"github.com/peteclark-io/footie/players"
)

var acknowledgementTemplate = template.Must(template.New("ack-email").Parse(acknowledgementContent))

const acknowledgementContent = `
<style>
  @import url('https://fonts.googleapis.com/css?family=Orienta:400|Mukta:400,700|Yantramanav:700,400&display=swap');

  * {
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-rendering: optimizeLegibility;
  }

  body {
    color: #444;
    margin: 0;
  }

  .wrap {
    display: flex;
    flex-direction: row;
    justify-content: center;
    background-color: #ddd;
  }

  .content {
    width: 76vw;
    background-color: #eee;
    text-align: center;
    padding: 1rem;
    padding-bottom: 2rem;
  }

  h1 {
    font-family: 'Yantramanav', sans-serif;
    font-size: 32pt;
    font-weight: 700;
    margin: 2rem;
  }

  h2,h3,h4 {
    font-family: 'Orienta', sans-serif;
    font-weight: 400;
    text-transform: uppercase;
    font-size: 16pt;
  }

  p {
    font-family: 'Yantramanav', sans-serif;
    font-size: 12pt;
    font-weight: 400;
  }

  a.btn {
    display: inline-block;
    padding: 0.7em 1.4em;
    margin: 0 0.5em 0.5em 0;
    border-radius: 0.15em;
    box-sizing: border-box;
    text-decoration: none;
    font-family: 'Yantramanav',sans-serif;
    text-transform: uppercase;
    font-weight: 400;
    color: #FFFFFF;
    background-color: #3369ff;
    box-shadow: inset 0 -0.6em 0 -0.35em rgba(0,0,0,0.17);
    text-align: center;
    position: relative;
    min-width: 100px;
  }

  a.btn-ok {
    background-color: #37ab2e;
  }

  a.btn-critical {
    background-color: #ff5252;
  }

  a.btn:active {
    top: 0.1em;
  }

  .details p {
    font-family: 'Yantramanav', sans-serif;
    font-size: 8pt;
    text-align: right;
  }

  @media all and (max-width:30em){
    h1 {
      font-size: 28pt;
    }
    a.btn {
      /* display:block; */
      min-width: 100px;
      margin: 0.4em 0.4em;
    }
  }
</style>
<div class="wrap">
  <section class="content">
    <h1><b>Monday Night Football</b></h1>
    <h3>Are you in for {{ .Match.Name }}?</h3>
    <p>To reply, click the appropriate button below, and send the email it creates.</p>
    <p style="display: none">When is it? 09-12-2019 at 18:20 to 19:00</p>
    <section>
      <a class="btn btn-ok" href="mailto:in+{{ .Match.Group }}{{ .Match.ID }}{{ .Player.ID }}@shoreditch.football?subject=IN&body=IN">IN</a>
      <a class="btn btn-critical" href="mailto:in+{{ .Match.Group }}{{ .Match.ID }}{{ .Player.ID }}@shoreditch.football?subject=OUT&body=OUT">OUT</a>
    </section>
  </section>
</div>
`

const CharSet = "UTF-8"
const DefaultSender = `"Monday Night Football" <in@shoreditch.football>`

func GenerateAcknowledgementEmail(player *players.Player, match *matches.Match) (string, error) {
	if player == nil || match == nil {
		return "", errors.New("Invalid input, please ensure neither player nor match are nil")
	}

	emailData := bytes.NewBufferString("")
	params := make(map[string]interface{})
	params["Match"] = match
	params["Player"] = player

	err := acknowledgementTemplate.Execute(emailData, params)
	if err != nil {
		return "", err
	}

	return string(emailData.Bytes()), nil
}

type EmailProducer struct {
}

func NewEmailProducer() *EmailProducer {
	return &EmailProducer{}
}

func (e *EmailProducer) ProduceEmail(sender, recipient, subject, body, textBody string) error {
	sess := session.Must(session.NewSession())

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(body),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	_, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return err
	}

	return nil
}
