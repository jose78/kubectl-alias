package database

import (
	"fmt"
	"strings"

	"github.com/jose78/kubectl-alias/commons"
	"github.com/jose78/sqlparser"
)


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
		commons.ErrorSqlNotASelect.BuildMsgError(query).KO()
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


// Find the tables used within the from
func FindTablesWithAliases(query string) map[string]string {
	query = strings.ReplaceAll(query, ".", "____")
	// Parsear la consulta a un AST
	stmt, _ := sqlparser.Parse(query)
	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		commons.ErrorSqlNotASelect.BuildMsgError(query).KO()
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