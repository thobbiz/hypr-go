/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package workspace

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/thobbiz/hypr-go/hyprerr"
	"github.com/thobbiz/hypr-go/workspace"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save <name> [workspace id]",
	Short: "Saves a workspace layout into a named profile",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		dir, err := getDir(cmd)
		if err != nil {
			return err
		}

		var workspaceID *int
		if len(args) == 2 {
			workspaceID, err = parseWorkspaceID(args[1])
			if err != nil {
				return err
			}
		}
		return workspace.SaveWorkspace(workspaceID, name, dir)
	},
}

func parseWorkspaceID(arg string) (*int, error) {
	n, err := strconv.Atoi(arg)
	if err != nil {
		return nil, hyprerr.New("parseWorkspaceID", hyprerr.ErrInvalidInput, fmt.Errorf("invalid workspace id %q: %w", arg, err))
	}
	return &n, nil
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
	// saveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// saveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
