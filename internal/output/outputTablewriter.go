package output

import (
	"fmt"
	"os"

	"github.com/jose78/kubectl-fuck/internal/database"
	"github.com/olekukonko/tablewriter"
)

func PrintStdout(result database.SelectResult) {

	columns := append([]string{"Index"}, result.Columns...)
	data := [][]string{}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columns)

	for index, rowMap := range result.Rows {
		row := []string{fmt.Sprintf("%d", index)}
		for _, col := range result.Columns {
			row = append(row, rowMap[col].(string))
		}
		data = append(data, row)
	}
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
