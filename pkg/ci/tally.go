package ci

import (
	"fmt"

	"github.com/alexeyco/simpletable"
)

func TallyTable(tally []ComponentTally) string {
	var table = simpletable.New()
	tallyHeader(table)
	tallyRows(tally, table)
	fmt.Println(table.String())

	return table.String()
}

func tallyHeader(table *simpletable.Table) {
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Component Type"},
			{Align: simpletable.AlignCenter, Text: "Count"},
		},
	}
}

func tallyRows(tally []ComponentTally, table *simpletable.Table) {

	for _, t := range tally {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: t.ComponentType},
			{Text: fmt.Sprintf("%v", t.Count)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

}
