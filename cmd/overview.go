package cmd

import (
	"errors"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"innodb_inspector/pkg/innodb"
	"innodb_inspector/pkg/innodb/page"
)

func overView(cmd *cobra.Command, filePath string) error {
	pds, err := innodb.OverView(filePath, int(page.DefaultSize))
	if err != nil {
		return errors.New("bad innodb file")
	}

	cells := make([][]*simpletable.Cell, 0)

	for _, pd := range pds {
		cells = append(cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", pd.PageNo)},
			{Text: fmt.Sprintf("%s", pd.PageType)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", pd.SpaceId)},
			{Text: pd.PageNotes},
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
			{Align: simpletable.AlignCenter, Text: "PageNotes"},
		},
	}

	table.Body.Cells = cells

	cmd.Println(table.String())
}
