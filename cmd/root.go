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
		conf := dn.ReadConf()
		path := conf.Path
		if path == "" {
			var err error
			path, err = cmd.Flags().GetString("path-to-note")
			if err != nil {
				panic(err)
			}
		}
		dn.Enter(path, conf.Editor)
	},
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
}
