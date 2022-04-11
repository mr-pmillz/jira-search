/*
Copyright © 2022 mr-pmillz

*/
package main

import (
	"log"
	"os"

	"github.com/mr-pmillz/jira-search/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	if val, present := os.LookupEnv("GENERATE_JIRA_SEARCH_DOCS"); val == "true" && present {
		if err := doc.GenMarkdownTree(cmd.RootCmd, "./docs"); err != nil {
			log.Fatal(err)
		}
	}
	cmd.Execute()
}
