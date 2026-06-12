/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package layout

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/thobbiz/hypr-go/layout"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save <name> <workspaceNumber>",
	Short: "Save specified workspace layout",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		var id *int
		if len(args) == 2 {
			n, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			id = &n
		}
		return layout.SaveWorkspace(id, name)
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// saveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
