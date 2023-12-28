package cmd

import (
	"context"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/kasaikou/thedocuments/tools"
	"github.com/kasaikou/thedocuments/tools/readme/core"
	"github.com/m-mizutani/clog"
	"github.com/spf13/cobra"
)

var rootLogger = slog.New(clog.New(
	clog.WithColor(true),
	clog.WithSource(true),
))

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
		logger := rootLogger.With(slog.String("rootCommand", "build"))
		ctx := tools.WithLogger(context.Background(), logger)
		wd, _ := os.Getwd()
		if err := core.Build(ctx, filepath.Join(wd, "readme.config.yaml")); err != nil {
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
