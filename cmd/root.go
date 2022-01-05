package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var filePath = ""
var pageNo int32 = -1
var exportPath = ""

var rootCmd = &cobra.Command{
	Use:   "innodb_inspector -f <file> [-p <page>]",
	Short: "InnoDB offline file browsing utility",
	RunE: func(cmd *cobra.Command, args []string) error {
		if filePath == "" {
			return errors.New("option '-f' requires an argument")
		}

		if pageNo >= 0 {
			pageDetail(cmd, filePath, uint32(pageNo), exportPath)
			return nil
		}

		return overView(cmd, filePath)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "file path")
	rootCmd.Flags().Int32VarP(&pageNo, "page", "p", -1, "page no")
	rootCmd.Flags().StringVar(&exportPath, "export", "", "export page binary")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
