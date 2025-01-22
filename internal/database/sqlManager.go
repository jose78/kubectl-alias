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
package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jose78/kubectl-alias/commons"
	"github.com/jose78/sqlparser"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type SelectResult struct {
	Columns []string
	Rows    []map[string]any
}

// Execute Select winthin the database
func EvaluateSelect(db *sql.DB, sqlSelect string) SelectResult {

	rows, errSelect := db.Query(sqlSelect)
	if errSelect != nil {
		commons.ErrorSqlRuningSelect.BuildMsgError(sqlSelect, errSelect).KO()
	}
	defer rows.Close() // Ensure rows are closed even if errors occur
	columns, err := rows.Columns()
	if err != nil {
		commons.ErrorSqlReadingColumns.BuildMsgError(err).KO()
	}
	results := SelectResult{Columns: columns, Rows: []map[string]any{}}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			commons.ErrorSqlScaningResultSelect.BuildMsgError(err).KO()
		}
		row := map[string]any{}
		for i, col := range columns {
			var v interface{}
			if b, ok := values[i].([]byte); ok {
				v = string(b) // Convertir []byte a string
			} else {
				v = values[i]
			}
			row[col] = v
		}
		results.Rows = append(results.Rows, row)
	}
	return results
}

// CReate table
func CreateTable(db *sql.DB, table string) {

	data := `CREATE TABLE %s (
        id INTEGER PRIMARY KEY,
        %s TEXT NOT NULL
    );
`
	createTable := fmt.Sprintf(data, table, table)
	statement, err := db.Prepare(createTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Printf("%s table created\n", table)
}

// Insert list of items of same type in a table
func Insert(db *sql.DB, k8sValues []unstructured.Unstructured, tbl string) {

	for _, value := range k8sValues {

		valueJson, _ := json.Marshal(value)
		valueStr := fmt.Sprintf("INSERT INTO %s(%s) VALUES('%s');", tbl, tbl, string(valueJson))
		statement, err := db.Prepare(valueStr) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			log.Fatalln(err.Error())
		}
		_, err = statement.Exec()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	log.Printf("Added into table %s %d elements", tbl, len(k8sValues))

}

func FindTablesWithAliases(query string) map[string]string {
	query = strings.ReplaceAll(query, ".", "____")
	// Parsear la consulta a un AST
	stmt, _ := sqlparser.Parse(query)
	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		log.Fatalf("El tipo de consulta no es un SELECT, es: %T", stmt)
	}

	aliasToTable := map[string]string{}
	for _, tableExpr := range selectStmt.From {
		aliasedTable, ok := tableExpr.(*sqlparser.AliasedTableExpr)
		if !ok {
			continue
		}

		tableName := aliasedTable.Expr.(sqlparser.TableName).Name.String()

		alias := tableName
		if aliasedTable.As.String() != "" {
			alias = aliasedTable.As.String()
		}

		aliasToTable[alias] = tableName
	}
	return aliasToTable
}

type colInfo struct {
	columnName string
	tableName  string
}

func regenerateColInfo(col string, aliasToTable map[string]string) colInfo {

	colSplited := strings.Split(col, "____")
	qualifier := ""
	name := ""
	if len(aliasToTable) != 1 {
		qualifier = aliasToTable[colSplited[0]]
		name = strings.Join(colSplited[1:], ".")
	} else {
		for _, value := range aliasToTable {
			qualifier = value
			break
		}
		_, ok := aliasToTable[colSplited[0]]
		if ok {
			name = strings.Join(colSplited[1:], ".")
		} else {
			name = strings.Join(colSplited, ".")
		}
	}

	tableName := qualifier
	columnName := fmt.Sprintf("$.Object.%s", name)

	return colInfo{columnName, tableName}
}

func UpdateQuery(query string, aliasToTable map[string]string) string {

	query = strings.ReplaceAll(query, ".", "____")
	// Parsear la consulta a un AST
	stmt, _ := sqlparser.Parse(query)

	updateSelect := func(aliasToTable map[string]string, selectStmt *sqlparser.Select) {
		for i, expr := range selectStmt.SelectExprs {
			aliasedExpr, ok := expr.(*sqlparser.AliasedExpr)
			if !ok {
				continue
			}

			colName, ok := aliasedExpr.Expr.(*sqlparser.ColName)
			if !ok {
				continue
			}

			col := regenerateColInfo(colName.Name.String(), aliasToTable)
			tableName := col.tableName
			columnName := col.columnName

			aliasedExpr.Expr = &sqlparser.FuncExpr{
				Name: sqlparser.NewColIdent("json_extract"),
				Exprs: sqlparser.SelectExprs{
					&sqlparser.AliasedExpr{
						Expr: &sqlparser.ColName{Name: sqlparser.NewColIdent(tableName)},
					},
					&sqlparser.AliasedExpr{
						Expr: &sqlparser.ColName{Name: sqlparser.NewColIdent(columnName)},
					},
				},
			}

			selectStmt.SelectExprs[i] = aliasedExpr
		}
	}

	updateWhere := func(aliasToTable map[string]string, selectStmt *sqlparser.Select) {
		if selectStmt.Where == nil {
			return
		}

		updaCol := func(binaryExpr *sqlparser.Expr) {

			colName, _ := (*binaryExpr).(*sqlparser.ColName)

			col := regenerateColInfo(colName.Name.String(), aliasToTable)
			tableName := col.tableName
			columnName := col.columnName

			*binaryExpr = &sqlparser.FuncExpr{
				Name: sqlparser.NewColIdent("json_extract"),
				Exprs: sqlparser.SelectExprs{
					&sqlparser.AliasedExpr{
						Expr: &sqlparser.ColName{Name: sqlparser.NewColIdent(tableName)},
					},
					&sqlparser.AliasedExpr{
						Expr: &sqlparser.ColName{Name: sqlparser.NewColIdent(columnName)},
					},
				},
			}
		}
		// Loop through the conditions of the WHERE clause
		sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
			binaryExpr, ok := node.(*sqlparser.ComparisonExpr)
			if !ok {
				return true, nil
			}

			_, ok = binaryExpr.Left.(*sqlparser.ColName)
			if ok {
				updaCol(&binaryExpr.Left)
			}

			_, ok = binaryExpr.Right.(*sqlparser.ColName)
			if ok {
				updaCol(&binaryExpr.Right)
			}

			return true, nil
		}, selectStmt.Where.Expr)
	}

	updateGroupBy := func(aliasToTable map[string]string, selectStmt *sqlparser.Select) {
		for i, expr := range selectStmt.GroupBy {
			colName, ok := expr.(*sqlparser.ColName)
			if !ok {
				continue
			}

			col := regenerateColInfo(colName.Name.String(), aliasToTable)
			tableName := col.tableName
			columnName := col.columnName

			selectStmt.GroupBy[i] = &sqlparser.FuncExpr{
				Name: sqlparser.NewColIdent("json_extract"),
				Exprs: sqlparser.SelectExprs{
					&sqlparser.AliasedExpr{
						Expr: &sqlparser.ColName{Name: sqlparser.NewColIdent(tableName)},
					},
					&sqlparser.AliasedExpr{
						Expr: &sqlparser.ColName{Name: sqlparser.NewColIdent(columnName)},
					},
				},
			}
		}
	}

	updateOrderBy := func(aliasToTable map[string]string, selectStmt *sqlparser.Select) {
		for i, order := range selectStmt.OrderBy {
			colName, ok := order.Expr.(*sqlparser.ColName)
			if !ok {
				continue
			}

			col := regenerateColInfo(colName.Name.String(), aliasToTable)
			tableName := col.tableName
			columnName := col.columnName

			selectStmt.OrderBy[i].Expr = &sqlparser.FuncExpr{
				Name: sqlparser.NewColIdent("json_extract"),
				Exprs: sqlparser.SelectExprs{
					&sqlparser.AliasedExpr{
						Expr: &sqlparser.ColName{Name: sqlparser.NewColIdent(tableName)},
					},
					&sqlparser.AliasedExpr{
						Expr: &sqlparser.ColName{Name: sqlparser.NewColIdent(columnName)},
					},
				},
			}
		}
	}

	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		log.Fatalf("El tipo de consulta no es un SELECT, es: %T", stmt)
	}

	updateSelect(aliasToTable, selectStmt)
	updateWhere(aliasToTable, selectStmt)
	updateGroupBy(aliasToTable, selectStmt)
	updateOrderBy(aliasToTable, selectStmt)
	modifiedQuery := sqlparser.String(stmt)
	modifiedQuery = strings.ReplaceAll(modifiedQuery, "`", "'")

	keywords := []string{
		"select ", "from ", "where ", "join ", " on ", "group by ", "order by ", "having ",
		"set", "limit", "offset",
	}

	// Reemplazar cada palabra clave por su versión en mayúsculas
	for _, keyword := range keywords {
		upper := strings.ToUpper(keyword)
		modifiedQuery = strings.ReplaceAll(modifiedQuery, keyword, upper)
	}

	return modifiedQuery
}
