package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"github.com/xwb1989/sqlparser"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type row map[string] any 
type selectResult struct{
	columns []string
	rows map[int]row
}

// Execute SElect
func evaluateSelect(db *sql.DB, sqlSelect string) (selectResult) {

	rows, errSelect := db.Query(sqlSelect)
	if errSelect != nil {
		ErrorSqlRuningSelect.buildMsgError(sqlSelect, errSelect).KO()
	}
	defer rows.Close() // Ensure rows are closed even if errors occur
	columns, err := rows.Columns()
	if err != nil {
		ErrorSqlReadingColumns.buildMsgError(err).KO()
	}
	results :=  selectResult{columns: columns, rows: map[int]row{}}
	index := 0
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			ErrorSqlScaningResultSelect.buildMsgError(err).KO()
		}
		row := row{}
		for i, col := range columns {
			var v interface{}
			if b, ok := values[i].([]byte); ok {
				v = string(b) // Convertir []byte a string
			} else {
				v = values[i]
			}
			row[col] = v
		}
		results.rows[index] = row
		index++
	}
	return results
}

// CReate table
func createTable(db *sql.DB, table string) {

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

// insert list of items of same type in a table
func insert(db *sql.DB, k8sValues []unstructured.Unstructured, tbl string) {

	for _, value := range k8sValues {

		valueJson, _ := json.Marshal(value)
		valueStr := fmt.Sprintf("INSERT INTO %s(%s) VALUES('%s');", tbl,tbl,  string(valueJson))
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

}
func evaluateQuery(sqlStr string) ([]string, error) {
	var evaluateFrom func(map[string]any) []string
	evaluateFrom = func(data map[string]any) []string {
		result := []string{}
		for key, value := range data {
			fmt.Println("key: ", key)
			if key == "Expr"{
				table := value.(map[string]any)["Name"].(string)
				result = append(result, table)
			} else if value != nil &&  reflect.TypeOf(value).Kind()  == reflect.Map {
				result = append(result, evaluateFrom(value.(map[string]any))...)
			}
		}
		return result
	}
	stmt, err := sqlparser.Parse(sqlStr)
	if err != nil {
		panic(err)
	}
	bytes, _ := json.Marshal(stmt)
	var data map[string]any
	json.Unmarshal(bytes, &data)

	result := []string{}
	from := data["From"]
	lstFrom := from.([]any)
	for _, itemFrom := range lstFrom {
		result = append(result, evaluateFrom(itemFrom.(map[string]any))...)
	}

	return result, nil
}
