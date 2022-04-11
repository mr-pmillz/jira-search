package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mr-pmillz/jira-search/cmd/search"

	"github.com/mr-pmillz/jira-search/utils"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	defaultConfigFileName = "config"
	envPrefix             = "JIRA-SEARCH"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "jira-search",
	Short: "Search Jira Issues From the Command Line",
	Long: `Search Jira Issues via jql queries and others from
the cli.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.yaml)")

	RootCmd.AddCommand(search.RootSearchCommand)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		absConfigFilePath, err := utils.ResolveAbsPath(cfgFile)
		if err != nil {
			_ = fmt.Errorf("Couldn't resolve path of config file: %v\n", err)
			return
		}
		viper.SetConfigFile(absConfigFilePath)
	} else {
		// Search config in project root directory with name "config.yaml" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultConfigFileName)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if _, err = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed()); err != nil {
			return
		}
	}

	viper.SetEnvPrefix(envPrefix)

	viper.AutomaticEnv() // read in environment variables that match

	bindFlags(RootCmd)

}

// bindFlags Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
		err := viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		if err != nil {
			return
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				return
			}
		}
	})
}
