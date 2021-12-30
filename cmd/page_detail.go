package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"innodb_inspector/pkg/innodb"
	"innodb_inspector/pkg/page"
)

func pageDetail(cmd *cobra.Command, filePath string, pageNo uint32) error {
	detail, err := innodb.PageDetail(filePath, pageNo, page.DefaultSize)
	if err != nil {
		return errors.New("bad innodb file")
	}

	cmd.Println(detail)
	return nil
}
