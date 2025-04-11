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
	"os"
	"path"
	"path/filepath"

	"github.com/jose78/kubectl-alias/commons"
	"github.com/jose78/kubectl-alias/internal/utils"
	_ "github.com/mattn/go-sqlite3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type SelectResult struct {
	Columns []string
	Rows    []map[string]any
}

// Execute Select winthin the database
func (conf dbConf) EvaluateSelect(sqlSelect string) SelectResult {

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
	utils.Logger(utils.INFO, "Create table " + table)
	data := `CREATE TABLE %s (
        id INTEGER PRIMARY KEY,
        %s TEXT NOT NULL
    );
`
	createTable := fmt.Sprintf(data, table, table)
	statement, err := conf.db.Prepare(createTable)
	if err != nil {
		commons.ErrorDBCreateTable.BuildMsgError(table, err)
	}
	statement.Exec()
}

// Insert list of items of same type in a table
func (conf dbConf) Insert(k8sValues []unstructured.Unstructured, table string) {
	elements := "elements"
	if len(k8sValues) == 1 {
		elements = "element"
	}
	utils.Logger(utils.INFO, fmt.Sprintf( "Insert %d %s in table %s" ,  len(k8sValues), elements ,table))
	for _, value := range k8sValues {

		valueJson, _ := json.Marshal(value)
		valueStr := fmt.Sprintf("INSERT INTO %s(%s) VALUES('%s');", table, table, string(valueJson))
		statement, err := conf.db.Prepare(valueStr) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			commons.ErrorDBInsertPrepare.BuildMsgError(table, err)
		}
		_, err = statement.Exec()
		if err != nil {
			commons.ErrorDBRunningInsert.BuildMsgError(table, err)

		}
	}
}

// colInfo represents column information, including its name and the table it belongs to.
type colInfo struct {
	columnName string
	tableName  string
}

type dbConf struct {
	db *sql.DB
}
type DbConf interface {
	CreateTable(string)
	Insert([]unstructured.Unstructured, string)
	EvaluateSelect(string) SelectResult
	Destroy()
}

func getDbFaile() string {
	kubeAliasPaths := os.Getenv(commons.ENV_VAR_KUBEALIAS_NAME)
	path := path.Dir(kubeAliasPaths)
	path = filepath.Join(path, "sqlite-database.db")
	return path
}

func Load() DbConf {
	utils.Logger(utils.INFO, "Load db")
	path := getDbFaile()

	file, errCReateDbObj := os.Create(path) // Create SQLite file
	if errCReateDbObj != nil {
		commons.ErrorDbNotCreaterd.BuildMsgError(errCReateDbObj).KO()
	}
	file.Close()

	sqliteDatabase, errOpeningDB := sql.Open("sqlite3", path) // Open the created SQLite File
	if errOpeningDB != nil {
		commons.ErrorDbOpening.BuildMsgError(errOpeningDB).KO()
	}
	conf := dbConf{db: sqliteDatabase}
	return conf
}

func (conf dbConf) Destroy() {
	conf.db.Close()
	os.Remove(getDbFaile())
	utils.Logger(utils.INFO, "removed database file: " +  getDbFaile())
}
