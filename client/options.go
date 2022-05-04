package client

import (
	"github.com/mr-pmillz/jira-search/utils"
	"github.com/spf13/cobra"
)

type Options struct {
	JiraHost        string
	JiraUserEmail   string
	JiraAPIKey      string
	JiraAccountID   string
	JiraUserName    string
	JiraProjectName string
	JQLTextSearch   string
	JQLRawSearch    string
	MyJiraIssues    bool
	RerverseIssues  bool
}

// ConfigureCommand ...
func ConfigureCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("jira-host", "", "", "jira host url ex. https://foobar.atlassian.net")
	cmd.PersistentFlags().StringP("jira-user-email", "", "", "jira user email address")
	cmd.PersistentFlags().StringP("jira-api-key", "", "", "jira api key")
	cmd.PersistentFlags().StringP("jira-account-id", "", "", "jira account id of api key user")
	cmd.PersistentFlags().StringP("jira-username", "", "", "jira username")
	cmd.PersistentFlags().StringP("jira-project-name", "", "", "a project name in jira you want to search")
	cmd.PersistentFlags().StringP("jql-text-search", "", "", "a string of text that you want to search all issues for")
	cmd.PersistentFlags().StringP("jql-raw-search", "", "", "a string of text that you want to search all issues for")
	cmd.PersistentFlags().BoolP("my-jira-issues", "", false, "find all issues assigned and in progress assigned to JIRA_ACCOUNT_ID from config.yaml")
	cmd.PersistentFlags().BoolP("reverse-issues", "", false, "reverse found issues in table")
}

// LoadFromCommand loads all the command flag opts from cli and config file into Options struct
func (opts *Options) LoadFromCommand(cmd *cobra.Command) error {
	myJiraIssues, err := cmd.Flags().GetBool("my-jira-issues")
	if err != nil {
		return err
	}
	opts.MyJiraIssues = myJiraIssues

	reverseIssues, err := cmd.Flags().GetBool("reverse-issues")
	if err != nil {
		return err
	}
	opts.RerverseIssues = reverseIssues

	jiraHost, err := utils.ConfigureFlagOpts(cmd, &utils.LoadFromCommandOpts{
		Flag: "jira-host",
		Opts: opts.JiraHost,
	})
	if err != nil {
		return err
	}
	opts.JiraHost = jiraHost.(string)

	jiraUserEmail, err := utils.ConfigureFlagOpts(cmd, &utils.LoadFromCommandOpts{
		Flag: "jira-user-email",
		Opts: opts.JiraUserEmail,
	})
	if err != nil {
		return err
	}
	opts.JiraUserEmail = jiraUserEmail.(string)

	jiraAPIKey, err := utils.ConfigureFlagOpts(cmd, &utils.LoadFromCommandOpts{
		Flag: "jira-api-key",
		Opts: opts.JiraAPIKey,
	})
	if err != nil {
		return err
	}
	opts.JiraAPIKey = jiraAPIKey.(string)

	jiraAccountID, err := utils.ConfigureFlagOpts(cmd, &utils.LoadFromCommandOpts{
		Flag:       "jira-account-id",
		IsFilePath: false,
		Prefix:     "",
		Opts:       opts.JiraAccountID,
	})
	if err != nil {
		return err
	}
	opts.JiraAccountID = jiraAccountID.(string)

	jiraUserName, err := utils.ConfigureFlagOpts(cmd, &utils.LoadFromCommandOpts{
		Flag:       "jira-username",
		IsFilePath: false,
		Prefix:     "",
		Opts:       opts.JiraUserName,
	})
	if err != nil {
		return err
	}
	opts.JiraUserName = jiraUserName.(string)

	jiraProjectName, err := utils.ConfigureFlagOpts(cmd, &utils.LoadFromCommandOpts{
		Flag:       "jira-project-name",
		IsFilePath: false,
		Prefix:     "",
		Opts:       opts.JiraProjectName,
	})
	if err != nil {
		return err
	}
	opts.JiraProjectName = jiraProjectName.(string)

	jqlTextSearch, err := utils.ConfigureFlagOpts(cmd, &utils.LoadFromCommandOpts{
		Flag:       "jql-text-search",
		IsFilePath: false,
		Prefix:     "",
		Opts:       opts.JQLTextSearch,
	})
	if err != nil {
		return err
	}
	opts.JQLTextSearch = jqlTextSearch.(string)

	jqlRawSearch, err := utils.ConfigureFlagOpts(cmd, &utils.LoadFromCommandOpts{
		Flag:       "jql-raw-search",
		IsFilePath: false,
		Prefix:     "",
		Opts:       opts.JQLRawSearch,
	})
	if err != nil {
		return err
	}
	opts.JQLRawSearch = jqlRawSearch.(string)

	return nil
}
