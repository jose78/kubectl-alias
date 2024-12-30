package main

import (
	"context"
	"fmt"

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

	tables, _ := evaluateQuery(aliasFiltered[0].SQL)
	ctx = context.WithValue(ctx, CTE_TABLES, tables)
	dataExtractedK8s :=  retrieveK8sObjects(ctx)

	
	for table, jsonContent := range dataExtractedK8s{
		createTable(sqliteDatabase, table)
		insert(sqliteDatabase,jsonContent , table )
	}
	dataSelect_1 := evaluateSelect(sqliteDatabase, aliasFiltered[0].SQL)

	fmt.Println(dataSelect_1)
}
