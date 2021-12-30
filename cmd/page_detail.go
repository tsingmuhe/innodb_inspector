package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"innodb_inspector/innodb/tablespace"
)

func pageDetail(cmd *cobra.Command, filePath string, pageNo uint32) error {
	ts := tablespace.NewTablespace(filePath)
	detail, err := ts.PageDetail(pageNo)
	if err != nil {
		return errors.New("bad innodb file")
	}

	cmd.Println(string(detail))
	return nil
}
