package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sterrasi/stepwise/logging"
)

// ApplicationName name of the appmication
const ApplicationName = "stepwise"

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "stepwise",
	Short: "Starts the stepwise application.",
	Long: `Starts the stepwise application.
			If the application configuration file is not specified explicitly then
			it is searched for in the following locations:
				1. /etc/stepwise/stepwise.toml
				2. $HOME/.stepwise/stepwise.toml
				3. ./stepwise.toml`,

	// load the config file for all commands
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// set config name and search paths
			viper.SetConfigName(ApplicationName)
			viper.AddConfigPath(fmt.Sprintf("/etc/%s", ApplicationName))
			viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", ApplicationName))
			viper.AddConfigPath(".")
		}

		// set defaults
		viper.SetDefault("server.cert-cache-dir", "/var/www/.cache")
		viper.SetDefault("server.address", ":443")
		viper.SetDefault("users.default-results-per-page", "20")
		viper.SetDefault("logging.level", logging.InfoLogLevel)
		viper.SetDefault("logging.format", logging.TextLoggingFormat)
		viper.SetDefault("logging.log-requests", false)

		viper.AutomaticEnv() // read in environment variables that match

		err := viper.ReadInConfig()
		if err != nil {
			return fmt.Errorf("Unable to load config file: %s", err)
		}
		return nil
	},
}

func init() {

	// explicitly provide the application config through a persistant flag
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
}

// Execute the root command of the application
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
