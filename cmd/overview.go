package cmd

import (
	"errors"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"innodb_inspector/innodb/tablespace"
)

func overView(cmd *cobra.Command, filePath string) error {
	ts := tablespace.NewTablespace(filePath)
	rows, err := ts.OverView()
	if err != nil {
		return errors.New("bad innodb file")
	}

	cells := make([][]*simpletable.Cell, 0)

	for _, row := range rows {
		cells = append(cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.PageNo)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.FilPageOffset)},
			{Text: fmt.Sprintf("%s", row.FilPageType)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.SpaceId)},
			{Text: row.Comments},
		})
	}

	printTable(cmd, cells)
	return nil
}

func printTable(cmd *cobra.Command, cells [][]*simpletable.Cell) {
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "PageNo"},
			{Align: simpletable.AlignCenter, Text: "PageOffset"},
			{Align: simpletable.AlignCenter, Text: "PageType"},
			{Align: simpletable.AlignCenter, Text: "SpaceId"},
			{Align: simpletable.AlignCenter, Text: "Comments"},
		},
	}

	table.Body.Cells = cells

	cmd.Println(table.String())
}
