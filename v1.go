package main

import (
	"context"

	"github.com/jose78/go-collections"
)

type AliasV1 struct {
	Name string `yaml:"name"`
	SQL  string `yaml:"sql"`
}
type AliasDefV1 struct {
	Version string    `yaml:"version"`
	Aliases []AliasV1 `yaml:"aliases"`
}

func (alias AliasDefV1) execute(ctx context.Context) {
	aliasFiltered := []AliasV1{}
	aliasName := ctx.Value(CTX_KEY_ALIAS_NAME).(string)
	collections.Filter(func(item AliasV1) bool { return item.Name == aliasName }, alias.Aliases, &aliasFiltered)

	if len(aliasFiltered) == 0 {
		ErrorKubeAliasNotFoud.buildMsgError(aliasName).KO()
	} else if len(aliasFiltered) > 1 {
		ErrorKubeAliasDuplicated.buildMsgError(aliasName).KO()
	}

	aliasToTable := findTablesWithAliases(aliasFiltered[0].SQL)
	tables := []string{}
	collections.Map(func(touple collections.Touple) any { return touple.Value }, aliasToTable, &tables)
	for _, table := range tables {
		ctx = context.WithValue(ctx, CTE_TABLE, table)
		jsonContent := retrieveK8sObjects(ctx)
		createTable(sqliteDatabase, table)
		insert(sqliteDatabase, jsonContent, table)
	}
	sqlSelect := updateQuery(aliasFiltered[0].SQL, aliasToTable)
	dataSelect_1 := evaluateSelect(sqliteDatabase, sqlSelect)

	// fmt.Println(dataSelect_1)
	printStdout(dataSelect_1)
}
