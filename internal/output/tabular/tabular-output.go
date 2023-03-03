package tabular

import (
	"strconv"

	"github.com/carbonetes/diggity/internal/output/save"
	"github.com/carbonetes/diggity/internal/output/util"
	"github.com/carbonetes/diggity/internal/parser/bom"

	"github.com/alexeyco/simpletable"
)

// PrintTable Packages in Table format
func PrintTable() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "NAME"},
			{Align: simpletable.AlignCenter, Text: "VERSION"},
			{Align: simpletable.AlignCenter, Text: "TYPE"},
		},
	}

	var cells [][]*simpletable.Cell

	// Sort packages alphabetically
	util.SortPackages()

	for i, p := range bom.Packages {
		i++
		cells = append(cells, []*simpletable.Cell{
			{Text: strconv.Itoa(i)},
			{Text: p.Name},
			{Text: p.Version},
			{Text: p.Type},
		})
	}

	totalPackages := strconv.Itoa(len(bom.Packages))

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Span: 2},
		{Align: simpletable.AlignCenter, Text: "Total Packages"},
		{Align: simpletable.AlignCenter, Text: totalPackages},
	}}

	table.SetStyle(simpletable.StyleDefault)

	if len(*bom.Arguments.OutputFile) > 0 {
		save.ResultToFile(table.String())
	} else {
		table.Println()
	}
}
