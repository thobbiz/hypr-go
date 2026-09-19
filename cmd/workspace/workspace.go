/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package workspace

import (
	"fmt"

	"github.com/spf13/cobra"
)

// WorkspaceCmd represents the workspace command
var WorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manages workspace layouts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("workspace called")
	},
}

func init() {
	WorkspaceCmd.AddCommand(saveCmd)
	WorkspaceCmd.AddCommand(restoreCmd)

	WorkspaceCmd.PersistentFlags().String("dir", "", "directory for layout files (default: $XDG_DATA_HOME/hypr-go/workspace)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// layoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// layoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
