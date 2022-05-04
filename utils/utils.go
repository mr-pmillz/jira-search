package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type LoadFromCommandOpts struct {
	DefaultFlagVal string
	Flag           string
	IsFilePath     bool
	Prefix         string
	Opts           interface{}
}

// ConfigureFlagOpts sets the cobra flag option to the LoadFromCommandOpts.Opts key
// it returns the parsed value of the cobra flag from LoadFromCommandOpts.Flag
func ConfigureFlagOpts(cmd *cobra.Command, LCMOpts *LoadFromCommandOpts) (interface{}, error) {
	cmdFlag, err := cmd.Flags().GetString(fmt.Sprintf("%s%s", LCMOpts.Prefix, LCMOpts.Flag))
	if err != nil {
		return nil, err
	}

	switch cmdFlag {
	case "":
		configVal := viper.GetString(strings.ToUpper(strings.ReplaceAll(fmt.Sprintf("%s%s", LCMOpts.Prefix, LCMOpts.Flag), "-", "_")))
		envVal, ok := os.LookupEnv(configVal)
		if ok {
			if LCMOpts.IsFilePath {
				fileExists, err := Exists(envVal)
				if err != nil {
					return nil, err
				}
				if fileExists {
					absVal, err := ResolveAbsPath(envVal)
					if err != nil {
						return nil, err
					}
					LCMOpts.Opts = absVal
				} else {
					LCMOpts.Opts = envVal
				}
			} else {
				LCMOpts.Opts = envVal
			}
		} else {
			if configVal != "" {
				if LCMOpts.IsFilePath {
					absConfigVal, err := ResolveAbsPath(configVal)
					if err != nil {
						return nil, err
					}
					LCMOpts.Opts = absConfigVal
				} else {
					LCMOpts.Opts = configVal
				}
			} else {
				if LCMOpts.DefaultFlagVal != "" && LCMOpts.IsFilePath {
					absDefaultVal, err := ResolveAbsPath(LCMOpts.DefaultFlagVal)
					if err != nil {
						return nil, err
					}
					_, err = os.Stat(absDefaultVal)
					if os.IsNotExist(err) {
						LCMOpts.Opts = cmdFlag
					} else {
						LCMOpts.Opts = absDefaultVal
					}
				} else if LCMOpts.DefaultFlagVal != "" && !LCMOpts.IsFilePath {
					LCMOpts.Opts = LCMOpts.DefaultFlagVal
				} else {
					LCMOpts.Opts = cmdFlag
				}
			}
		}
	default:
		envValue, ok := os.LookupEnv(strings.ToUpper(strings.ReplaceAll(fmt.Sprintf("%s%s", LCMOpts.Prefix, LCMOpts.Flag), "-", "_")))
		if ok {
			LCMOpts.Opts = envValue
		} else {
			if LCMOpts.IsFilePath {
				fileExists, err := Exists(cmdFlag)
				if err != nil {
					return nil, err
				}
				if fileExists {
					absCmdFlag, err := ResolveAbsPath(cmdFlag)
					if err != nil {
						return nil, err
					}
					LCMOpts.Opts = absCmdFlag
				} else {
					LCMOpts.Opts = cmdFlag
				}

			} else {
				LCMOpts.Opts = cmdFlag
			}
		}
	}

	return LCMOpts.Opts, nil
}

// isDone ...
func isDone(status string) tablewriter.Colors {
	if status == "Done" || status == "Closed" {
		return tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}
	}
	return tablewriter.Colors{tablewriter.Bold, tablewriter.Normal}
}

// getPriorityColor ...
func getPriorityColor(priority string) tablewriter.Colors {
	switch priority {
	case "Highest":
		return tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}
	case "High":
		return tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor}
	case "Medium":
		return tablewriter.Colors{tablewriter.Bold, tablewriter.Normal}
	case "Low":
		return tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiMagentaColor}
	default:
		return tablewriter.Colors{tablewriter.Bold, tablewriter.Normal}
	}
}

func getFieldString(name interface{}) string {
	if name != nil {
		return name.(string)
	}
	return "N/A"
}

// PrintTable ...
func PrintTable(issues []jira.Issue, host string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Issue", "Summary", "Project", "Status", "Priority", "URL", "Reporter"})
	table.SetBorder(false)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor, tablewriter.BgBlackColor},
	)

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor, tablewriter.BgBlackColor},
	)

	for _, issue := range issues {
		colorData := []string{issue.Key,
			issue.Fields.Summary,
			issue.Fields.Project.Name,
			issue.Fields.Status.Name,
			issue.Fields.Priority.Name,
			fmt.Sprintf("%s/browse/%s", host, issue.Key),
			getFieldString(issue.Fields.Reporter.DisplayName),
		}

		table.Rich(colorData, []tablewriter.Colors{
			{tablewriter.Normal},
			{tablewriter.Normal},
			{tablewriter.Normal},
			isDone(issue.Fields.Status.Name),
			getPriorityColor(issue.Fields.Priority.Name),
			{tablewriter.Normal},
			{tablewriter.Normal},
		})
	}

	table.SetAutoMergeCells(true)
	table.Render()

	return nil
}

// ResolveAbsPath ...
func ResolveAbsPath(path string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return path, err
	}

	dir := usr.HomeDir
	if path == "~" {
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return path, err
	}

	return path, nil
}

// Exists returns whether the given file or directory exists
func Exists(path string) (bool, error) {
	absPath, err := ResolveAbsPath(path)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(absPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func WriteAttachmentBodyToLocalTmpFile(body io.ReadCloser, fileName string) (string, error) {
	fmt.Printf("[+] Downloading attachment: %s\n", fileName)
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err = io.Copy(tmpFile, body); err != nil {
		return "", err
	}

	tmpFilePath, err := ResolveAbsPath(tmpFile.Name())
	if err != nil {
		return "", err
	}

	return tmpFilePath, nil
}
