package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "get test code",
	Long:  `this is GET code`,
        Run: func(cmd *cobra.Command, args []string) {

        },


}

func init() {

	RootCmd.AddCommand(GetCmd)

}
