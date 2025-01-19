package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func printStdout(result selectResult) {

	columns := append([]string{"Index"}, result.columns...)
	data := [][]string{}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columns)

	for index, rowMap := range result.rows {
		row := []string{fmt.Sprintf("%d", index)}
		for _, col := range result.columns {
			row = append(row, rowMap[col].(string))
		}
		data = append(data, row)
	}
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
