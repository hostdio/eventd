package httpserver

import (
	"log"
	"net/http"

	"context"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/hostdio/eventd/api"
	"github.com/hostdio/eventd/plugins/database/inmemory"
	"github.com/hostdio/eventd/plugins/database/postgres"
	"github.com/hostdio/eventd/queues/googlepubsub"
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
	cmd.AddCommand(cmdLocal())

	return cmd
}

func cmdLocal() *cobra.Command {
	return &cobra.Command{
		Use:   "local",
		Short: "local starts the event sourcing HTTP API server",
		Long:  `local exposes the event sourcing through an HTTP API server. `,
		Run: func(cmd *cobra.Command, args []string) {
			publisher := inmemory.NewQueue()
			persistor := inmemory.NewDatabase()

			r := getmux(persistor, persistor, publisher, publisher)

			log.Fatal(http.ListenAndServe(":8080", r))
		},
	}
}

func cmdProduction() *cobra.Command {
	return &cobra.Command{
		Use:   "production <project ID> <topic ID>",
		Short: "production starts the event sourcing HTTP API server",
		Long:  `production exposes the event sourcing through an HTTP API server. `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 4 {
				fmt.Println("Missing either projectID, topicID, subID, or connStr")
				return
			}
			projectID := args[0]
			topicID := args[1]
			subID := args[2]
			connStr := args[3]

			publisher, err := googlepubsub.New(context.Background(), projectID, topicID, subID)
			if err != nil {
				panic(err)
			}
			defer publisher.Close()
			persistor, dbErr := postgres.New(connStr)
			if dbErr != nil {
				panic(dbErr)
			}

			r := getmux(persistor, persistor, publisher, publisher)

			log.Fatal(http.ListenAndServe(":8080", r))
		},
	}
}

func getmux(scanner api.Scanner, persister api.Persister, publisher api.Publisher, listener api.Listener) *mux.Router {
	r := mux.NewRouter()

	go startPersister(context.Background(), listener, persister)

	r.HandleFunc("/", publishHandler(publisher)).Methods("POST")
	r.HandleFunc("/", scanHandler(scanner)).Methods("GET")

	return r
}
