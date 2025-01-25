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
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jose78/go-collections"
	"github.com/jose78/kubectl-alias/internal/generic"
	"github.com/jose78/kubectl-alias/service"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubectl-alias",
	Short: "Customize your Kubernetes queries using SQL, including JOINs across different resources",
	Long: `This tool empowers you to customize your queries to Kubernetes resources using SQL. 
With support for JOIN operations, you can combine and filter data from multiple Kubernetes objects 
(Pods, ConfigMaps, Services, and more) seamlessly.

Key features:
- Write SQL queries to interact with Kubernetes resources.
- Perform JOINs across different Kubernetes object types.
- Customize subcommands and their behavior via a configuration file referenced by the KUBEALIAS environment variable.

Usage example:
kubectl alias subCommand arg1 arg2... argN

By defining subcommands (with or without arguments) in the KUBEALIAS configuration file, 
you can tailor the CLI to suit your specific Kubernetes query needs.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	contentKubeAlias := service.LoadKubeAlias()

	aliases, okAliases := contentKubeAlias["aliases"]
	if !okAliases {
		return
	}

	for name, value := range aliases.(map[string]any) {


		cmdCtx := generic.CommandContext{SubCommand: name}
		

		mapperArg := func(value string) any {
			return fmt.Sprintf("[%s]", strings.ToUpper(strings.ReplaceAll(value, " ", "_")))
		}

		item := value.(map[string]any)
		short := fmt.Sprintf("%s", item["short"])
		long := fmt.Sprintf("%s", item["long"])

		use := name
		sizeArgs := 0
		args, okArgs := item["args"]
		if okArgs {
			sizeArgs = len(args.([]string))
			var useList []string
			collections.Map(mapperArg, args.([]string), &useList)
			use = fmt.Sprintf("%s %s ", name, strings.Join(useList, " "))
		}

		var subCmd = &cobra.Command{
			Use:   use,
			Short: short,
			Long:  long,
			Args:  cobra.ExactArgs(sizeArgs), // Enforce exactly the len of the arguments
			Run: func(cmd *cobra.Command, args []string) {
				cmdCtx.Args = args
				service.RunAlias(cmdCtx)
			},
		}

		rootCmd.AddCommand(subCmd)

	}
}
