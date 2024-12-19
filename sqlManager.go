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

// Execute SElect
func evaluateSelect(db *sql.DB, sqlSelect string) ([][]interface{}, error) {

	rows, errSelect := db.Query(sqlSelect)
	if errSelect != nil {
		return nil, errSelect
	}
	defer rows.Close() // Ensure rows are closed even if errors occur

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %w", err)
	}

	var results [][]interface{} // Store results as a slice of slices (interface{})

	// Prepare destination variables for scanning
	dest := make([]interface{}, len(columnNames))
	for i, _ := range columnNames {
		dest[i] = new(interface{}) // Allocate memory for each interface{}
	}

	for rows.Next() {
		// Scan row values into destination variables
		if err := rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert scanned values to desired types if necessary
		rowValues := make([]interface{}, len(columnNames))
		for _, val := range dest {
			dataVal := *val.(*interface{})
			rowValues = append(rowValues, dataVal)
		}
		results = append(results, rowValues) // Append the row values to the results slice
	}
	return results, nil
}

// CReate table
func createTable(db *sql.DB, table string) {

	data := `CREATE TABLE %s(
        id INTEGER PRIMARY KEY,
        data TEXT NOT NULL
    );
`
	createTable := fmt.Sprintf(data, table)
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
		valueStr := fmt.Sprintf("INSERT INTO %s(data) VALUES('%s');", tbl, string(valueJson))
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
		for _, value := range data {
			if value != nil {
				kind := reflect.TypeOf(value).Kind()
				if kind == reflect.Map {
					if item, ok := value.(map[string]any)["Expr"]; ok {
						table := item.(map[string]any)["Name"].(string)
						result = append(result, table)
					} else if value != "Condition" {
						result = append(result, evaluateFrom(value.(map[string]any))...)
					}
				}
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
