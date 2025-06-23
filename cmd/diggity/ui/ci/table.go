package ci

import (
	"fmt"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/alexeyco/simpletable"
	"github.com/carbonetes/diggity/pkg/types"
)

func Table(cdx *cyclonedx.BOM) string {
	var table = simpletable.New()
	header(table)
	rows(cdx, table)
	fmt.Println(table.String())
	return table.String()
}

func header(table *simpletable.Table) {
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Package Name"},
			{Align: simpletable.AlignCenter, Text: "Type"},
			{Align: simpletable.AlignCenter, Text: "Version"},
		},
	}
}

func rows(cdx *cyclonedx.BOM, table *simpletable.Table) {
	for _, c := range *cdx.Components {
		var componentType string
		for _, p := range *c.Properties {
			if p.Name == "diggity:package:type" {
				if p.Value != "" {
					componentType = p.Value
				}
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

// Secrets Table
func SecretTable(secrets []types.Secret) string {
	var table = simpletable.New()
	secretHeader(table)
	secretRows(&secrets, table)
	fmt.Println(table.String())
	return table.String()
}

func secretHeader(table *simpletable.Table) {
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Secret"},
			{Align: simpletable.AlignCenter, Text: "File"},
			{Align: simpletable.AlignCenter, Text: "Line"},
			{Align: simpletable.AlignCenter, Text: "Content"},
			{Align: simpletable.AlignCenter, Text: "Type"},
		},
	}
}

func secretRows(secrets *[]types.Secret, table *simpletable.Table) {
	for _, s := range *secrets {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: s.Match},
			{Text: fmt.Sprintf("%v", s.File)},
			{Text: fmt.Sprintf("%v", s.Line)},
			{Text: fmt.Sprintf("%v", s.Content)},
			{Text: fmt.Sprintf("%v", s.Description)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
}
