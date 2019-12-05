package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/husobee/vestigo"
	"github.com/jawher/mow.cli"
	"github.com/peteclark-io/footie/acks"
	"github.com/peteclark-io/footie/aws"
	"github.com/peteclark-io/footie/bookings"
	"github.com/peteclark-io/footie/groups"
	"github.com/peteclark-io/footie/matches"
	"github.com/peteclark-io/footie/players"
	"github.com/peteclark-io/footie/resources"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := cli.App("footie", "Books football matches for players")

	resource := app.String(cli.StringOpt{
		Name:   "resource",
		Desc:   "The resource type to run",
		EnvVar: "RESOURCE",
	})

	method := app.String(cli.StringOpt{
		Name:   "method",
		Desc:   "The resource method to handle",
		EnvVar: "METHOD",
	})

	log.SetLevel(log.InfoLevel)

	app.Command("http", "Runs a webserver for local testing", func(cmd *cli.Cmd) {
		startHTTP(players.NewHTTPHandler(), acks.NewHTTPHandler(), matches.NewHTTPHandler(), groups.NewHTTPHandler(), bookings.NewHTTPHandler())
	})

	app.Action = func() {
		var handler aws.Handler
		switch *resource {
		case "players":
			handler = players.NewAWSHandler()
		case "matches":
			handler = matches.NewAWSHandler()
		case "groups":
			handler = groups.NewAWSHandler()
		case "bookings":
			handler = bookings.NewAWSHandler()
		case "acknowledgements":
			handler = acks.NewAWSHandler()
		}

		startAWS(*method, handler)
	}

	log.SetLevel(log.InfoLevel)
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("App could not start, error=[%s]\n", err)
		return
	}

	app.Run(os.Args)
}

func startHTTP(handlers ...resources.Handler) {
	r := vestigo.NewRouter()

	r.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin:  []string{"http://localhost:3000"},
		AllowHeaders: []string{"content-type"},
	})

	for _, h := range handlers {
		r.Get("/"+h.Name()+"/:id", h.Read)
		r.Post("/"+h.Name(), h.Create)
		r.Put("/"+h.Name()+"/:id", h.Write)
		r.Delete("/"+h.Name()+"/:id", h.Delete)
	}

	log.Fatal(http.ListenAndServe(":8000", resources.AccessLog(log.StandardLogger(), r)))
}

func startAWS(method string, handler aws.Handler) {
	switch method {
	case "create":
		lambda.Start(handler.Create)
	case "write":
		lambda.Start(handler.Write)
	case "read":
		lambda.Start(handler.Read)
	case "delete":
		lambda.Start(handler.Delete)
	}
}
