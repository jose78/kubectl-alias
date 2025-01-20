package alias

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/jose78/go-collections"
	"github.com/jose78/kubectl-fuck/commons"
	"github.com/jose78/kubectl-fuck/internal/database"
	"github.com/jose78/kubectl-fuck/internal/k8s"
	"github.com/jose78/kubectl-fuck/internal/output"
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
	Version string    `yaml:"version"`
	Aliases []AliasV1 `yaml:"aliases"`
}

// Implementation of interface Command for version V1 of alias functionality
func (alias AliasDefV1) Execute(ctx context.Context) {
	aliasFiltered := []AliasV1{}
	aliasName := ctx.Value(commons.CTX_KEY_ALIAS_NAME).(string)
	collections.Filter(func(item AliasV1) bool { return item.Name == aliasName }, alias.Aliases, &aliasFiltered)

	if len(aliasFiltered) == 0 {
		commons.ErrorKubeAliasNotFoud.BuildMsgError(aliasName).KO()
	} else if len(aliasFiltered) > 1 {
		commons.ErrorKubeAliasDuplicated.BuildMsgError(aliasName).KO()
	}

	aliasToTable := database.FindTablesWithAliases(aliasFiltered[0].SQL)
	tables := []string{}
	collections.Map(func(touple collections.Touple) any { return touple.Value }, aliasToTable, &tables)
	for _, table := range tables {
		ctx = context.WithValue(ctx, commons.CTE_TABLE, table)
		jsonContent := k8s.RetrieveK8sObjects(ctx)
		database.CreateTable(sqliteDatabase, table)
		database.Insert(sqliteDatabase, jsonContent, table)
	}
	sqlSelect := database.UpdateQuery(aliasFiltered[0].SQL, aliasToTable)
	dataSelect_1 := database.EvaluateSelect(sqliteDatabase, sqlSelect)

	// fmt.Println(dataSelect_1)
	output.PrintStdout(dataSelect_1)
}
