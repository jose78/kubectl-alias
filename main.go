package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	collection "github.com/jose78/go-collections"
	"github.com/xwb1989/sqlparser"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func main() {
	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer sqliteDatabase.Close()                                     // Defer Closing the database
	objRetrieved, _ := k8sGetElements(nil)
	for table, value := range objRetrieved {
		createTable(sqliteDatabase, table) // Create Database Tables
		insert(sqliteDatabase, value, table)
	}

	sqlSelect := `SELECT 
    json_extract(p.data, '$.metadata.name') AS pod_name,
    json_extract(d.data, '$.metadata.name') AS deploy_name
FROM 
    pod AS p
JOIN 
    deploy AS d ON json_extract(p.data, '$.metadata.labels.app') = json_extract(d.data, '$.metadata.name')
;`

	dataSelect_1, _ := evaluateSelect(sqliteDatabase, sqlSelect)
	fmt.Println(dataSelect_1)
	defer os.Remove("sqlite-database.db")
}


// Retrieve list of elements requered from sql
func k8sGetElements(elements []string) (result map[string][]any, err error) {
	return
}


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
func insert(db *sql.DB, elementToIterate []any, tbl string) {

	for _, value := range elementToIterate {

		valueByte, errK8sObj := json.Marshal(value)
		if errK8sObj != nil {
			log.Fatal(errK8sObj)
		}

		valueStr := fmt.Sprintf("INSERT INTO %s(data) VALUES('%s');", tbl, string(valueByte))
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

type sqlContainer struct {
	sqlSelect map[string]string
	sqlFrom   []string
	sqlWhere  string
}

func evaluateQuery(query string) (sqlContainer, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		log.Fatalf("Error al parsear la consulta SQL: %v", err)
	}

	sql := sqlContainer{}
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		cols := map[string]string{}
		for _, col := range stmt.SelectExprs {
			col := sqlparser.String(col)
			var key string
			var value string
			if upperNameItemSelect := strings.ToUpper(col); strings.Contains(upperNameItemSelect, " AS ") {
				key = strings.Split(upperNameItemSelect, " AS ")[0]
				value = strings.Split(upperNameItemSelect, " AS ")[1]
			} else {
				colSplited := strings.Split(col, ".")
				key = colSplited[len(colSplited)-1]
				value = col
			}
			cols[key] = value
			sql.sqlSelect = cols
		}
		sql.sqlFrom = strings.Split(sqlparser.String(stmt.From), ",")
		sql.sqlWhere = strings.ReplaceAll(sqlparser.String(stmt.Where), " where ", "")
	default:
		return sqlContainer{}, fmt.Errorf("")
	}
	return sql, nil
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
		return "", fmt.Errorf(`error: not found env var %s`, envNameVar)
	} else {
		return result[envNameVar], nil
	}
}
