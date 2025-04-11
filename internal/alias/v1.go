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
	"fmt"
	"strings"

	"github.com/jose78/go-collections"
	"github.com/jose78/kubectl-alias/commons"
	"github.com/jose78/kubectl-alias/internal/database"
	"github.com/jose78/kubectl-alias/internal/generic"
	"github.com/jose78/kubectl-alias/internal/k8s"
	"github.com/jose78/kubectl-alias/internal/output"
	"github.com/jose78/kubectl-alias/internal/utils"
	"github.com/spf13/cobra"
)

type AliasV1 struct {
	Short string   `yaml:"short"`
	Long  string   `yaml:"long"`
	Args  []string `yaml:"args"`
	SQL   string   `yaml:"sql"`
}

type AliasDefV1 struct {
	Version string             `yaml:"version"`
	Aliases map[string]AliasV1 `yaml:"aliases"`
}

// Implementation of interface Command for version V1 of alias functionality
func (aliasFiltered AliasV1) execute(ctx generic.CommandContext) {
	utils.Logger(utils.INFO, "Start v1")
	sql := aliasFiltered.SQL
	if len(aliasFiltered.Args) > 0 {
		for index := 0; index < len(aliasFiltered.Args); index++ {
			sql = strings.ReplaceAll(sql, aliasFiltered.Args[index], ctx.Args[index])
		}
	}

	aliasToTable := database.FindTablesWithAliases(sql)
	tables := []string{}
	collections.Map(func(touple collections.Touple) any { return touple.Value }, aliasToTable, &tables)
	pathK8sConfig := ctx.Flags[commons.CTE_KUBECONFIG].(*string)
	namespaceDefault := ctx.Flags[commons.CTE_NS].(*string)
	k8sInfo := k8s.K8sInfo{PathK8sConfig: *pathK8sConfig, NamespaceDefault: *namespaceDefault}
	mapObjects := k8s.GenerateMapObjects(k8sInfo)
	k8sInfo.K8sResources = mapObjects

	sqlSelect := database.ManipulateAST(sql, aliasToTable)

	dbObjetc := database.Load()
	defer dbObjetc.Destroy()
	for _, table := range tables {
		jsonContent := k8s.RetrieveK8sObjects(k8sInfo, table)
		dbObjetc.CreateTable(table)
		dbObjetc.Insert(jsonContent, table)
	}
	rows := dbObjetc.EvaluateSelect(sqlSelect)
	output.PrintStdout(rows)
	utils.Logger(utils.INFO, "End v1")
}

func (alias AliasDefV1) GenerateDoc(ctx generic.CommandContext) []*cobra.Command {

	result := []*cobra.Command{}

	for name, value := range alias.Aliases {
		mapperArg := func(value string) any {
			return fmt.Sprintf("[%s]", value)
		}
		use := name
		sizeArgs := 0

		if len(value.Args) > 0 {
			sizeArgs = len(value.Args)
			var useList []string
			collections.Map(mapperArg, value.Args, &useList)
			use = fmt.Sprintf("%s %s ", name, strings.Join(useList, " "))
		}
		var subCmd = &cobra.Command{
			Use:   use,
			Short: value.Short,
			Long:  value.Long,
			Args:  cobra.ExactArgs(sizeArgs),
			Run: func(cmd *cobra.Command, args []string) {
				aliasFiltered, okAlias := alias.Aliases[name]
				if !okAlias {
					commons.ErrorKubeAliasNotFoud.BuildMsgError(name).KO()
				}
				ctx.Args = args
				aliasFiltered.execute(ctx)
			},
		}
		result = append(result, subCmd)
	}

	return result
}
