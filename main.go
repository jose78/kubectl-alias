package main

import (
	"fmt"
	collection "github.com/jose78/go-collections"
	"github.com/xwb1989/sqlparser"
	"log"
	"os"
	"strings"
)

func main() {

}

func mapper(item string) any {
	stringSplited := strings.Split(item, "=")
	value := ""
	key := stringSplited[0]
	if len(stringSplited) > 1 {
		value = stringSplited[1]
	}
	return collection.Touple{Key: key, Value: value}
}

// function to extract the query to be executed
func selectQuery(nameVar string) (string, error) {
	envNameVar := fmt.Sprintf("K_FCK_%s", strings.ToUpper(nameVar))
	var predicate collection.Predicate[collection.Touple] = func(item collection.Touple) bool {
		result := strings.Compare(item.Key.(string), envNameVar) == 0
		return result
	}

	lstKeysEnv := os.Environ()
	mapUpdated := map[string]string{}
	collection.Map(mapper, lstKeysEnv, mapUpdated)
	result := map[string]string{}
	collection.Filter(predicate, mapUpdated, result)

	if len(result) == 0 {
		return "", fmt.Errorf(`Error: not found env var %s`, envNameVar)
	} else {
		return result[envNameVar], nil
	}
}

type sqlContainer struct {
	sqlSelect map[string]string
	sqlFrom []string 
	sqlWhere string 
}


func evaluateQuery(query string) (sqlContainer , error){
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		log.Fatalf("Error al parsear la consulta SQL: %v", err)
	}

	sql :=sqlContainer{}
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		cols := map[string]string{}
		for _, col := range stmt.SelectExprs {
			col := sqlparser.String(col)
			var key string
			var value string
			if upperNameItemSelect := strings.ToUpper(col);   strings.Contains(upperNameItemSelect, " AS ") {
				key   = strings.Split(upperNameItemSelect, " AS ")[0]
				value  = strings.Split(upperNameItemSelect, " AS ")[1]
			} else {
				colSplited := strings.Split(col, ".")
				key = colSplited[len(colSplited) - 1]
				value = col
			}
			cols[key] = value
			sql.sqlSelect = cols
		}
		sql.sqlFrom =  strings.Split(sqlparser.String(stmt.From) , ",")
		sql.sqlWhere = strings.ReplaceAll(sqlparser.String(stmt.Where ), " where ", "") 
	default:
		return sqlContainer{}, fmt.Errorf("")
	}
	return sql, nil
}
