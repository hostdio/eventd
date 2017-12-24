package httpserver

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/hostdio/eventd/queues/googlepubsub"
	"context"
	"fmt"
	"github.com/hostdio/eventd/databases/postgres"
)

// Cmd returns the command for starting an http server
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "httpserver <project ID> <topic ID>",
		Short: "httpserver starts the event sourcing HTTP API server",
		Long:  `httpserver exposes the event sourcing through an HTTP API server. `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 4 {
				fmt.Println("Missing either projectID, topicID, subID, or connStr")
				return
			}
			projectID := args[0]
			topicID := args[1]
			subID := args[2]
			connStr := args[3]



			r := mux.NewRouter()
			publisher, err := googlepubsub.New(context.Background(), projectID, topicID, subID)
			if err != nil {
				panic(err)
			}
			defer publisher.Close()
			persistor, dbErr := postgres.New(connStr)
			if dbErr != nil {
				panic(dbErr)
			}

			go startPersister(context.Background(), publisher, persistor)

			r.HandleFunc("/", publishHandler(publisher)).Methods("POST")
			log.Fatal(http.ListenAndServe(":8080", r))
		},
	}
}
