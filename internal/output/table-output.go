package output

import (
	"sort"
	"strconv"

	"github.com/carbonetes/diggity/internal/parser"

	"github.com/alexeyco/simpletable"
)

// Print Packages in Table format
func printTable() {
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
	sortPackages()

	for i, p := range parser.Packages {
		i++
		cells = append(cells, []*simpletable.Cell{
			{Text: strconv.Itoa(i)},
			{Text: p.Name},
			{Text: p.Version},
			{Text: p.Type},
		})
	}

	totalPackages := strconv.Itoa(len(parser.Packages))

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Span: 2},
		{Align: simpletable.AlignCenter, Text: "Total Packages"},
		{Align: simpletable.AlignCenter, Text: totalPackages},
	}}

	table.SetStyle(simpletable.StyleDefault)

	if len(*parser.Arguments.OutputFile) > 0 {
		saveResultToFile(table.String())
	} else {
		table.Println()
	}
}

// Sort packages by name alphabetically
func sortPackages() {
	sort.Slice(parser.Packages, func(i, j int) bool {
		return parser.Packages[i].Name < parser.Packages[j].Name
	})
}
