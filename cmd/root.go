package cmd

import (
	"github.com/ckut/legiosk/pkg/download"
	"github.com/spf13/cobra"
)

var lorientFlag string

func init() {
	RootCmd.PersistentFlags().StringVar(&lorientFlag, "lorient-username", "", "Login with your Lorient Mediatheque account")
}

var RootCmd = &cobra.Command{
	Use: "legiosk",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Help()
		}
		for _, url := range args {
			err := download.Download(url, lorientFlag)
			if err != nil {
				return err
			}
		}
		return nil
	},
}
