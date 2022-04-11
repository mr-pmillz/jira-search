package search

import (
	"fmt"
	"log"

	"github.com/mr-pmillz/jira-search/client"
	"github.com/spf13/cobra"
)

type SearchOptions = client.Options

type Options struct {
	SearchOptions
}

// LoadFromCommand ...
func (opts *Options) LoadFromCommand(cmd *cobra.Command) error {
	if err := opts.SearchOptions.LoadFromCommand(cmd); err != nil {
		return err
	}

	return nil
}

// RootSearchCommand ...
var RootSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search Jira for issues matching text",
	Long:  "Search Jira for issues matching text",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		opts := Options{}
		if err = opts.LoadFromCommand(cmd); err != nil {
			log.Panic(err)
		}

		jc, err := client.NewJiraClient(&opts.SearchOptions)
		if err != nil {
			log.Panic(err)
		}

		if err = jc.PrintFoundTickets(); err != nil {
			log.Panic(err)
		}

	},
}

// String ...
func (opts Options) String() string {
	return fmt.Sprintf("%v\n", opts.SearchOptions)
}

func init() {
	client.ConfigureCommand(RootSearchCommand)
}
