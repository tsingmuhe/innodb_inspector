package cmd

import (
	"errors"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"innodb_inspector/pkg/innodb"
)

func overView(cmd *cobra.Command, filePath string) error {
	ts := innodb.NewTablespace(filePath)
	pages, err := ts.OverView()
	if err != nil {
		return errors.New("bad innodb file")
	}

	cells := make([][]*simpletable.Cell, 0)

	for _, row := range pages {
		cells = append(cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.PageNo())},
			{Text: fmt.Sprintf("%s", row.Type())},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.SpaceId())},
			{Text: row.Notes()},
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
			{Align: simpletable.AlignCenter, Text: "PageType"},
			{Align: simpletable.AlignCenter, Text: "SpaceId"},
			{Align: simpletable.AlignCenter, Text: "Comments"},
		},
	}

	table.Body.Cells = cells

	cmd.Println(table.String())
}
