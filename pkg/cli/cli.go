package cli

import (
	"context"

	"github.com/cupsadarius/gcp_resource_cleaner/pkg/errors"
	"github.com/spf13/cobra"
)

var cmd *cobra.Command

// CommandHandlerFunc describes the header of functions that can be attached to a command
// All the functions passed to AddCommand must respect it
type CommandHandlerFunc func(ctx context.Context)

// Init initializes the CLI service
func Init(appID, shortDesc, longDesc string) {
	cmd = &cobra.Command{
		Use:   appID,
		Short: shortDesc,
		Long:  longDesc,
	}
}

// AddCommand adds a new command to CLI service
func AddCommand(command, description string, handlerFunc CommandHandlerFunc) error {
	if cmd == nil {
		return errors.ErrNotInitialized
	}

	cmd.AddCommand(&cobra.Command{
		Use:   command,
		Short: description,
		Long:  description,
		Run: func(cmd *cobra.Command, _ []string) {
			handlerFunc(cmd.Context())
		},
	})

	return nil
}

// AssignStringFlag set a string flag to CLI service
func AssignStringFlag(target *string, name, defaultValue, description string) {
	cmd.PersistentFlags().StringVar(target, name, defaultValue, description)
}

// AssignBoolFlag set a bool flag to CLI service
func AssignBoolFlag(target *bool, name string, defaultValue bool, description string) {
	cmd.PersistentFlags().BoolVar(target, name, defaultValue, description)
}

// Run runs the CLI service with a context attached
func Run(ctx context.Context) error {
	return cmd.ExecuteContext(ctx)
}
