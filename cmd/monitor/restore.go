/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package monitor

import (
	"github.com/spf13/cobra"
	"github.com/thobbiz/hypr-go/hyprerr"
	"github.com/thobbiz/hypr-go/monitor"
)

// monitor/restoreCmd represents the monitor/restore command
var restoreCmd = &cobra.Command{
	Use:   "restore <name>",
	Short: "Restores named monitor layout",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		dir, err := getDir(cmd)
		if err != nil {
			return err
		}

		return monitor.RestoreMonitor(name, dir)
	},
}

func getDir(cmd *cobra.Command) (string, error) {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return "", hyprerr.New("getDir", hyprerr.ErrParse, err)
	}
	return dir, nil
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitor/restoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitor/restoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
