package database

import (
	"fmt"
	"strings"

	"github.com/jose78/kubectl-alias/commons"
	"github.com/jose78/sqlparser"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const separatorLength = 8

// generateSeparator genera una cadena aleatoria de caracteres para usar como separador.
// Utiliza math/rand para la generación de números aleatorios.  Aunque rand.Seed está obsoleto,
// math/rand sigue siendo aceptable para este caso de uso donde no se requiere seguridad
// criptográfica.  Se usa un nuevo RandomSource para mayor control.
// La longitud del separador se define mediante la constante separatorLength.
func generateSeparator() string {
	// Se crea una nueva fuente de aleatoriedad con la hora actual como semilla.
	// Esto asegura que las cadenas generadas sean diferentes en cada ejecución.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, separatorLength)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))] // Usar el nuevo RandomSource
	}
	return "__________" + string(b) + "__________"
}

// ManipulateAST modifies the SQL query by transforming column references
// to use the `json_extract` function in SQLite3, using the table name as the key.
//
// SQLite3's `json_extract` function is used to extract values from JSON columns efficiently.
// It allows querying nested JSON fields directly within SQL queries.
//
// Parameters:
//   - query: The original SQL query string to be modified.
//   - aliasToTable: A map that associates table aliases with their actual table names.
//     If a table has no alias, its name is used as the alias.
//
// Returns:
//   - A modified SQL query string with column references wrapped in `json_extract`.
func ManipulateAST(query string, aliasToTable map[string]string) string {
	separator := generateSeparator()
	separatorConcat := fmt.Sprintf(" + '__________%s__________' + ", generateSeparator())
	replacer := strings.NewReplacer(".", separator ,"[", "open_______" ,"]", "close_______" , "||" , separatorConcat)
	query = replacer.Replace(query)
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return ""
	}

	visit := func(node sqlparser.SQLNode) (kontinue bool, err error) {
		if col, ok := node.(*sqlparser.ColName); ok {

			colInfo := regenerateColInfo(strings.Split(col.Name.String(), separator), aliasToTable)
			tableName := colInfo.tableName
			columnName := colInfo.columnName

			// Modificar el nombre de la columna para usar hisn()
			col.Name = sqlparser.NewColIdent(
				fmt.Sprintf("%sjson_extract(%s , '%s'%s)", separator, tableName, columnName, separator),
			)
		}
		return true, nil
	}
	sqlparser.Walk(visit, stmt)

	modifiedQuery := sqlparser.String(stmt)

	replacer = strings.NewReplacer("`"+separator, "" , separator+")`", ")"  , "open_______" , "[" , "close_______" , "]", separatorConcat, " || ")
	                                                                                                                                      

	fmt.Print(modifiedQuery)
	modifiedQuery = replacer.Replace(modifiedQuery)
	fmt.Print(modifiedQuery)
	return modifiedQuery
}

// regenerateColInfo reconstructs column information from a JSON path and table alias mapping.
//
// Parameters:
//   - colSplited: A slice containing the JSON path split by dots, representing the column hierarchy.
//   - aliasToTable: A map that associates table aliases with their actual table names.
//     If a table has no alias, its name is used as the alias.
//
// Returns:
//
//	A colInfo struct containing the column name and the corresponding table name.
func regenerateColInfo(colSplited []string, aliasToTable map[string]string) colInfo {

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

// FindTablesWithAliases extracts table names and their corresponding aliases from a SQL query.
//
// It analyzes the FROM clause, storing table aliases as keys in the returned map.
// If a table has no alias, its actual name is used as the key.
//
// Parameters:
//   - query: The SQL query string to analyze.
//
// Returns:
//   - A map where the key is the table alias (or the table name if no alias is used),
//     and the value is the actual table name.
func FindTablesWithAliases(query string) map[string]string {

	separator := generateSeparator()
	separatorConcat := fmt.Sprintf(" + '__________%s__________' + ", separator)


	replacer := strings.NewReplacer(".", separator ,"[", "open_______" ,"]", "close_______" , "||" , separatorConcat)


	fmt.Println(query)
	query = replacer.Replace(query)
	fmt.Println(query)
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
