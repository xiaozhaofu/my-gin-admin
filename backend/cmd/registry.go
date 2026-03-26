package cmd

import "github.com/spf13/cobra"

type commandFactory func(*rootOptions) *cobra.Command
type commandSetFactory func(*rootOptions) []*cobra.Command

var (
	rootCommandFactories []commandFactory
	dbCommandFactories   []commandSetFactory
)

// RegisterRootCommand 注册 root 级别命令扩展。
func RegisterRootCommand(factory commandFactory) {
	rootCommandFactories = append(rootCommandFactories, factory)
}

// RegisterDBCommands 注册 db 命令组下的一组子命令扩展。
func RegisterDBCommands(factory commandSetFactory) {
	dbCommandFactories = append(dbCommandFactories, factory)
}

func buildRootCommands(opts *rootOptions) []*cobra.Command {
	commands := make([]*cobra.Command, 0, len(rootCommandFactories))
	for _, factory := range rootCommandFactories {
		if factory == nil {
			continue
		}

		commands = append(commands, factory(opts))
	}

	return commands
}

func buildDBCommands(opts *rootOptions) []*cobra.Command {
	var commands []*cobra.Command
	for _, factory := range dbCommandFactories {
		if factory == nil {
			continue
		}

		commands = append(commands, factory(opts)...)
	}

	return commands
}
