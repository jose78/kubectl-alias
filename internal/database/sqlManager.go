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
	"os"

	"github.com/jose78/kubectl-alias/commons"
	_ "github.com/mattn/go-sqlite3" 
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type SelectResult struct {
	Columns []string
	Rows    []map[string]any
}

// Execute Select winthin the database
func (conf dbConf)EvaluateSelect(sqlSelect string) SelectResult {

	rows, errSelect := conf.db.Query(sqlSelect)
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
		for index, column := range columns {
			var v interface{}
			if b, ok := values[index].([]byte); ok {
				v = string(b)
			} else {
				v = values[index]
			}
			row[column] = v
		}
		results.Rows = append(results.Rows, row)
	}
	return results
}

// CReate table
func (conf dbConf) CreateTable(table string) {

	data := `CREATE TABLE %s (
        id INTEGER PRIMARY KEY,
        %s TEXT NOT NULL
    );
`
	createTable := fmt.Sprintf(data, table, table)
	statement, err := conf.db.Prepare(createTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Printf("%s table created\n", table)
}

// Insert list of items of same type in a table
func (conf dbConf) Insert( k8sValues []unstructured.Unstructured, tbl string) {

	for _, value := range k8sValues {

		valueJson, _ := json.Marshal(value)
		valueStr := fmt.Sprintf("INSERT INTO %s(%s) VALUES('%s');", tbl, tbl, string(valueJson))
		statement, err := conf.db.Prepare(valueStr) // Prepare statement.
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


type colInfo struct {
	columnName string
	tableName  string
}




type dbConf struct {
	db *sql.DB
}
type DbConf interface{
	CreateTable(string)
	Insert([]unstructured.Unstructured, string)
	EvaluateSelect(string) SelectResult
	Destroy()
} 



func Load() DbConf{
	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	file, errCReateDbObj := os.Create("sqlite-database.db") // Create SQLite file
	if errCReateDbObj != nil {
		commons.ErrorDbNotCreaterd.BuildMsgError(errCReateDbObj).KO()
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, errOpeningDB := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	if errOpeningDB != nil{
		commons.ErrorDbOpening.BuildMsgError(errOpeningDB).KO()
	}
	conf := dbConf{db: sqliteDatabase}
	return conf 
}

func (conf dbConf) Destroy() {
	conf.db.Close()
	os.Remove("./sqlite-database.db")
}
