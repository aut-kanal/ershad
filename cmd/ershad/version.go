package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/kanalbot/ershad"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Ershad's version",
	Run: func(cmd *cobra.Command, args []string) {
		logVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func logVersion() {
	logrus.Info("version   > ", ershad.Version)
	logrus.Info("buildtime > ", ershad.BuildTime)
	logrus.Info("commit    > ", ershad.Commit)
}
