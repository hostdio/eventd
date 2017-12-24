package httpserver

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/hostdio/eventd/queues/googlepubsub"
	"context"
	"fmt"
)

// Cmd returns the command for starting an http server
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "httpserver <project ID> <topic ID>",
		Short: "httpserver starts the event sourcing HTTP API server",
		Long:  `httpserver exposes the event sourcing through an HTTP API server. `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				fmt.Println("Missing either projectID or topicID")
				return
			}
			projectID := args[0]
			topicID := args[1]
			r := mux.NewRouter()
			publisher, err := googlepubsub.New(context.Background(), projectID, topicID)
			if err != nil {
				panic(err)
			}
			r.HandleFunc("/", publishHandler(publisher)).Methods("POST")
			log.Fatal(http.ListenAndServe(":8080", r))
		},
	}
}
