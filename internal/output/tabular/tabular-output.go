package tabular

import (
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/carbonetes/diggity/internal/output/save"
	"github.com/carbonetes/diggity/pkg/model"
)

// PrintTable Packages in Table format
func PrintTable(args *model.Arguments, outputType *string, pkgs *[]model.Package) {
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

	for i, p := range *pkgs {
		i++
		cells = append(cells, []*simpletable.Cell{
			{Text: strconv.Itoa(i)},
			{Text: p.Name},
			{Text: p.Version},
			{Text: p.Type},
		})
	}

	totalPackages := strconv.Itoa(len(*pkgs))

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Span: 2},
		{Align: simpletable.AlignCenter, Text: "Total Packages"},
		{Align: simpletable.AlignCenter, Text: totalPackages},
	}}

	table.SetStyle(simpletable.StyleCompactLite)

	if len(*args.OutputFile) > 0 {
		save.ResultToFile(table.String(), outputType, args.OutputFile)
	} else {
		table.Println()
	}
}
