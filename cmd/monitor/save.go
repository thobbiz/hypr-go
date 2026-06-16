/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package monitor

import (
	"github.com/spf13/cobra"
	"github.com/thobbiz/hypr-go/monitor"
)

// monitor/saveCmd represents the monitor/save command
var saveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "captures all currently connected monitors into a named profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return monitor.SaveMonitor(name)
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitor/saveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitor/saveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
