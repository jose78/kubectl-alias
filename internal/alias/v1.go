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
	"github.com/jose78/go-collections"
	"github.com/jose78/kubectl-alias/commons"
	"github.com/jose78/kubectl-alias/internal/database"
	"github.com/jose78/kubectl-alias/internal/generic"
	"github.com/jose78/kubectl-alias/internal/k8s"
	"github.com/jose78/kubectl-alias/internal/output"
)


type AliasV1 struct {
	Name string `yaml:"name"`
	SQL  string `yaml:"sql"`
}

type AliasDefV1 struct {
	Version string             `yaml:"version"`
	Aliases map[string]AliasV1 `yaml:"aliases"`
}

// Implementation of interface Command for version V1 of alias functionality
func (alias AliasDefV1) Execute(ctx generic.CommandContext) {
	aliasName := ctx.SubCommand
	aliasFiltered, okAlias := alias.Aliases[aliasName]
	if !okAlias {
		commons.ErrorKubeAliasNotFoud.BuildMsgError(aliasName).KO()
	}
	aliasToTable := database.FindTablesWithAliases(aliasFiltered.SQL)
	tables := []string{}
	collections.Map(func(touple collections.Touple) any { return touple.Value }, aliasToTable, &tables)
	
	k8sInfo := k8s.K8sInfo{ PathK8sConfig: ctx.Flags[commons.CTE_KUBECONFIG],  NamespaceDefault: ctx.Flags[commons.CTE_NS]}
	mapObjects := k8s.GenerateMapObjects(k8sInfo)
	k8sInfo.K8sResources = mapObjects

	sqlSelect := database.UpdateQuery(aliasFiltered.SQL, aliasToTable)
	
	dbObjetc := database.Load()
	defer dbObjetc.Destroy()
	for _, table := range tables {
		jsonContent := k8s.RetrieveK8sObjects(k8sInfo , table)
		dbObjetc.CreateTable(table)
		dbObjetc.Insert(jsonContent, table)
	}
	rows := dbObjetc.EvaluateSelect( sqlSelect)
	output.PrintStdout(rows)
}
