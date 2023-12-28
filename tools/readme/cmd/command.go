package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	buildCmd.PersistentFlags().StringVarP(&flagConfigFilename, "config", "c", "", "config file")
	rootCmd.AddCommand(buildCmd)
}

var (
	flagConfigFilename string
)

var rootCmd = &cobra.Command{
	Use: "readme",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("readme")
	},
}

var buildCmd = &cobra.Command{
	Use: "build",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("readme build")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
