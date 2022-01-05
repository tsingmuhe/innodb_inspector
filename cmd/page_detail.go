package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"innodb_inspector/pkg/innodb"
	"innodb_inspector/pkg/innodb/page"
)

func pageDetail(cmd *cobra.Command, filePath string, pageNo uint32, exportPath string) error {
	detail, err := innodb.PageDetail(filePath, pageNo, page.DefaultSize, exportPath)
	if err != nil {
		return errors.New("bad innodb file")
	}

	cmd.Println(detail)
	return nil
}
