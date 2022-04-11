package client

import (
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	//"github.com/dghubble/oauth1"
	"github.com/mr-pmillz/jira-search/utils"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type JiraClient struct {
	ctx    context.Context
	client *jira.Client
	opts   *Options
}

// NewJiraClient Creates wrapper for the Jira Client and everything else
func NewJiraClient(opts *Options) (*JiraClient, error) {
	j := &JiraClient{}

	// Create transport using the PersonalAccessToken specified in the password.
	tp := jira.BasicAuthTransport{
		Username: opts.JiraUserEmail,
		Password: opts.JiraAPIKey,
	}

	jiraClient, err := jira.NewClient(tp.Client(), opts.JiraHost)
	if err != nil {
		return nil, err
	}

	j.client = jiraClient

	if j.ctx == nil {
		j.ctx = context.Background()
	}

	j.opts = opts

	return j, nil
}

// PrintFoundTickets searches for matching ticket types and prints them
func (jc *JiraClient) PrintFoundTickets() error {
	// Ensure that there is a matching ticket for the current project with this JQL query.
	var jql string
	if jc.opts.JQLRawSearch != "" {
		jql = fmt.Sprintf("%s", jc.opts.JQLRawSearch)
	} else if jc.opts.MyJiraIssues {
		jql = fmt.Sprintf("status in (\"Ready for work\", \"In Progress\", \"Deploy Ready\") AND assignee in ( %s ) ORDER BY created ASC", jc.opts.JiraAccountID)
	} else {
		jql = fmt.Sprintf("text ~ \"%s\" ORDER BY created ASC", jc.opts.JQLTextSearch)
	}
	//jql := fmt.Sprintf("text ~ \"%s\" ORDER BY created ASC", jc.opts.JQLTextSearch)
	issues, resp, err := jc.client.Issue.Search(jql, &jira.SearchOptions{
		Fields: []string{
			"attachment",
			"key",
			"issuetype",
			"summary",
			"reporter",
			"status",
			"description",
			"project",
			"priority"},
	},
	)
	if err != nil {
		return err
	}

	total := resp.Total

	if total >= 1 {
		if err = utils.PrintTable(issues, jc.opts.JiraHost); err != nil {
			return nil
		}
	} else {
		fmt.Printf("[-] No Issues Found Matching this JQL Search: \n%s\n", jql)
		return nil
	}
	return nil
}

func (jc *JiraClient) DownloadAttachments(issue jira.Issue) (map[string]string, error) {

	patchFiles := make(map[string]string)

	attachments := issue.Fields.Attachments
	for _, attachment := range attachments {
		resp, err := jc.client.Issue.DownloadAttachment(attachment.ID)
		if err != nil {
			return nil, err
		}

		patchFile, err := utils.WriteAttachmentBodyToLocalTmpFile(resp.Body, attachment.Filename)
		if err != nil {
			return nil, err
		}
		err = resp.Body.Close()
		if err != nil {
			return nil, err
		}
		patchFiles[attachment.Filename] = patchFile
	}

	return patchFiles, nil
}

func (jc *JiraClient) UploadAttachment(issue jira.Issue, attachmentPath, attachmentName string) error {
	data, err := os.Open(attachmentPath)
	if err != nil {
		return err
	}
	reader := io.Reader(data)

	fmt.Printf("[+] Uploading attachment %s to %s\n", attachmentName, issue.Key)
	attachment, resp, err := jc.client.Issue.PostAttachment(issue.ID, reader, attachmentName)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Debug("error closing post attachment body")
		}
	}(resp.Body)

	if attachment == nil {
		logrus.Debug("Expected response. Response is nil")
	}

	return nil
}

// SetTicketStatus sets the ticket status to In Progress else quit.
func (jc *JiraClient) SetTicketStatus(issue jira.Issue, status string) error {
	issueKey := issue.Key
	fmt.Printf("[+] Setting %s Status to: %s\n", issueKey, status)
	var transitionID string
	possibleTransitions, _, _ := jc.client.Issue.GetTransitions(issueKey)
	for _, v := range possibleTransitions {
		if v.Name == status {
			transitionID = v.ID
			break
		}
	}

	_, err := jc.client.Issue.DoTransition(issueKey, transitionID)
	if err != nil {
		return err
	}
	return nil
}

func (jc *JiraClient) AddCommentToTicket(issueKey, commentBody string) error {
	comment := jira.Comment{}
	comment.Body = commentBody

	_, _, err := jc.client.Issue.AddComment(issueKey, &comment)
	if err != nil {
		return err
	}

	return nil
}

// AssignTicket Assign ticket to user
func (jc *JiraClient) AssignTicket(issue jira.Issue, userID string) error {
	issueKey := issue.Key

	jUser, _, err := jc.client.User.Get(userID)
	if err != nil {
		return err
	}

	fmt.Printf("[+] Assigning %s to %s\n", issueKey, jUser.DisplayName)
	_, err = jc.client.Issue.UpdateAssignee(issueKey, jUser)
	if err != nil {
		return err
	}

	return nil
}
