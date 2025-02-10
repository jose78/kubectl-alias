/*
Copyright © 2025 Jose Clavero Anderica (jose.clavero.anderica@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package output

import (
	"fmt"
	"os"

	"github.com/jose78/kubectl-alias/internal/database"
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
			row = append(row, fmt.Sprintf("%v", rowMap[col]))
		}
		data = append(data, row)
	}
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
