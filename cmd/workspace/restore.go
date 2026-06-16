/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package workspace

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/thobbiz/hypr-go/workspace"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore <name> <workspace id>",
	Short: "Restores specified workspace layout",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		var workspaceID *int
		if len(args) == 2 {
			n, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			workspaceID = &n
		}
		return workspace.RestoreWorkspace(name, workspaceID)
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
