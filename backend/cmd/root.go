package cmd

import (
	"fmt"

	"go_sleep_admin/internal/appmeta"

	"github.com/spf13/cobra"
)

type rootOptions struct {
	envFile string
}

func Execute() error {
	return NewRootCmd().Execute()
}

func NewRootCmd() *cobra.Command {
	opts := &rootOptions{}

	rootCmd := &cobra.Command{
		Use:           appmeta.CommandName,
		Short:         appmeta.RuntimeName + " runtime",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.PersistentFlags().StringVarP(&opts.envFile, "config", "c", "dev", "config environment name")
	rootCmd.AddCommand(buildRootCommands(opts)...)

	return rootCmd
}

func commandError(action string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", action, err)
}
