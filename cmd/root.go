package cmd

import (
	"os"

	"acakp.dn/dn"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nd",
	Short: "Create daily note",
	// Long: ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		conf := getCompleteConf(cmd)
		dn.Enter(conf)
	},
}

func getFlags(cmd *cobra.Command) dn.Config {
	var flags dn.Config
	flags.Path = getFlag(cmd, "path-to-note")
	flags.Editor = getFlag(cmd, "editor")
	flags.Extension = getFlag(cmd, "extension")
	flags.Format = getFlag(cmd, "format")
	flags.IsYesterday = getFlag(cmd, "yesterday")
	return flags
}

func getFlag(cmd *cobra.Command, flagName string) string {
	flag, err := cmd.Flags().GetString(flagName)
	if err != nil {
		panic(err)
	}
	return flag
}

func getCompleteConf(cmd *cobra.Command) dn.Config {
	flags := getFlags(cmd)
	conf := dn.ReadConf()
	if conf.Path == "" || cmd.Flags().Changed("path-to-note") {
		conf.Path = flags.Path
	}
	if conf.Editor == "" || cmd.Flags().Changed("editor") {
		conf.Editor = flags.Editor
	}
	if conf.Format == "" || cmd.Flags().Changed("format") {
		conf.Format = flags.Format
	}
	if conf.Extension == "" || cmd.Flags().Changed("extension") {
		conf.Extension = flags.Extension
	}
	return conf
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.acakp.nd.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("path-to-note", "p", "~/", "Path where note will be saved")
	rootCmd.Flags().StringP("editor", "e", "", "Editor to write a note")
	rootCmd.Flags().StringP("format", "f", "%YYYY-%M-%D %w", "How the note filename will look like")
	rootCmd.Flags().StringP("extension", "E", "md", "File extension")
	rootCmd.Flags().StringP("yesterday", "y", "false", "Open yesterday's note")
}
