package cmd

import (
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "eventd",
	Short: "eventd is an API for event sourcing",
	Long:  `eventd is a event sourcing API built on top of commodity infrastructure.`,
}

// Execute executes the cmd
func Execute() error {
	return root.Execute()
}
