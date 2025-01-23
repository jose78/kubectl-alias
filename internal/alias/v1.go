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

package alias

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/jose78/go-collections"
	"github.com/jose78/kubectl-alias/commons"
	"github.com/jose78/kubectl-alias/internal/database"
	"github.com/jose78/kubectl-alias/internal/k8s"
	"github.com/jose78/kubectl-alias/internal/output"
)

var (
	sqliteDatabase *sql.DB
)

func init() {
	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, _ = sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
}

func destroy() {
	sqliteDatabase.Close()
	os.Remove("./sqlite-database.db")
}

type AliasV1 struct {
	Name string `yaml:"name"`
	SQL  string `yaml:"sql"`
}
type AliasDefV1 struct {
	Version string             `yaml:"version"`
	Aliases map[string]AliasV1 `yaml:"aliases"`
}

// Implementation of interface Command for version V1 of alias functionality
func (alias AliasDefV1) Execute(ctx context.Context) {
	aliasName := ctx.Value(commons.CTX_KEY_ALIAS_NAME).(string)
	aliasFiltered, okAlias := alias.Aliases[aliasName]
	if !okAlias {
		commons.ErrorKubeAliasNotFoud.BuildMsgError(aliasName).KO()
	}
	aliasToTable := database.FindTablesWithAliases(aliasFiltered.SQL)
	tables := []string{}
	collections.Map(func(touple collections.Touple) any { return touple.Value }, aliasToTable, &tables)


	k8s.K8sConf{}

	mapObjects := k8s.GenerateMapObjects(ctx)
	ctx = context.WithValue(ctx, commons.CTE_MAP_K8S_OBJECT, mapObjects)

	for _, table := range tables {
		ctx = context.WithValue(ctx, commons.CTE_TABLE, table)
		jsonContent := k8s.RetrieveK8sObjects(ctx)
		database.CreateTable(sqliteDatabase, table)
		database.Insert(sqliteDatabase, jsonContent, table)
	}
	sqlSelect := database.UpdateQuery(aliasFiltered.SQL, aliasToTable)
	dataSelect_1 := database.EvaluateSelect(sqliteDatabase, sqlSelect)

	// fmt.Println(dataSelect_1)
	output.PrintStdout(dataSelect_1)
	defer destroy()
}
