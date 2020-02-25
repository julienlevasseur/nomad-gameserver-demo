package cmd

import (
	"github.com/julienlevasseur/nomad-gameserver-demo/matchmaker/matcmaker"
	"github.com/spf13/cobra"
)

// RootCmd root command
var RootCmd = &cobra.Command{
	Use:   "matchmaker",
	Short: "A simple POC matcmaker tool.",
	Run: func(cmd *cobra.Command, args []string) {
		matcmaker.CreateGameServer()
	},
}

// Execute is used in main.go
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}
