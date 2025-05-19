package ci

import (
	"fmt"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/alexeyco/simpletable"
)

func ResultTable(cdx *cyclonedx.BOM) string {

	var table = simpletable.New()
	resultHeader(table)
	resultRows(cdx, table)
	fmt.Println(table.String())

	return table.String()
}

func resultHeader(table *simpletable.Table) {
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Package Name"},
			{Align: simpletable.AlignCenter, Text: "Type"},
			{Align: simpletable.AlignCenter, Text: "Version"},
		},
	}
}

func resultRows(cdx *cyclonedx.BOM, table *simpletable.Table) {
	for _, c := range *cdx.Components {
		var componentType string
		for _, p := range *c.Properties {
			if p.Name == "diggity:package:type" {
				componentType = p.Value
			}
		}
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: c.Name},
			{Text: fmt.Sprintf("%v", componentType)},
			{Text: fmt.Sprintf("%v", c.Version)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
}
