package httpserver

import (
	"log"
	"net/http"

	"context"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/hostdio/eventd/api"
	"github.com/hostdio/eventd/eventkit"
	"github.com/hostdio/eventd/plugins/database/postgres"
	"github.com/spf13/cobra"
)

// Cmd returns the command for starting an http server
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "httpserver <project ID> <topic ID>",
		Short: "httpserver starts the event sourcing HTTP API server",
		Long:  `httpserver exposes the event sourcing through an HTTP API server. `,
	}

	cmd.AddCommand(cmdProduction())

	return cmd
}

func cmdProduction() *cobra.Command {
	return &cobra.Command{
		Use:   "production <project ID> <sub id> <conn str>",
		Short: "production starts the event sourcing HTTP API server",
		Long:  `production exposes the event sourcing through an HTTP API server. `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 4 {
				fmt.Println("Missing either projectID, topicID, subID, or connStr")
				return
			}
			projectID := args[0]
			subID := args[1]
			connStr := args[2]

			persistor, dbErr := postgres.New(connStr)
			if dbErr != nil {
				panic(dbErr)
			}

			consumer := eventkit.NewPubsubConsumer(context.Background(), projectID, subID)

			r := getmux(persistor, consumer, persistor)

			log.Fatal(http.ListenAndServe(":8080", r))
		},
	}
}

func getmux(scanner api.Scanner, consumer eventkit.Consumer, persister api.Persister) *mux.Router {
	r := mux.NewRouter()

	go startPersister(context.Background(), consumer, persister)

	r.HandleFunc("/", scanHandler(scanner)).Methods("GET")

	return r
}
