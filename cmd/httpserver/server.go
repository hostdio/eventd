package httpserver

import (
	"github.com/spf13/cobra"
	"log"
)

// Cmd returns the command for starting an http server
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "httpserver",
		Short: "httpserver starts the event sourcing HTTP API server",
		Long: `httpserver exposes the event sourcing through an HTTP API server. `,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Hello World")
		},
	}
}