package cmd

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/kasaikou/thedocuments/tools/readme/core"
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
		wd, _ := os.Getwd()
		if err := core.Build(context.Background(), filepath.Join(wd, "readme.config.yaml")); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
