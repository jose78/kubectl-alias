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

package commons

import (
	"fmt"
	"os"

	"github.com/jose78/kubectl-alias/internal/utils"
)

type errorManager int

type ErrorSystem struct {
	errorCode int
	errorMsg  string
}

const (
	ErrorK8sObjectnotSupported errorManager = iota
	ErrorK8sGeneratingDynamicClient
	ErrorK8sRestConfig
	ErrorK8sRestResource
	ErrorK8sRestResourceWithoutNS
	ErrorK8sClientConfig
	ErrorK8sNotContext
	ErrorK8sKubeconfgNotAccesible
	ErrorJsonMarshallResourceList
	ErrorKubeAliasPathNotDefined
	ErrorKubeAliasParseFile
	ErrorKubeAliasReadingFile
	ErrorKubeAliasVersionNotFoud
	ErrorKubeAliasNotFoud
	ErrorSqlRuningSelect
	ErrorSqlNotASelect
	ErrorSqlReadingColumns
	ErrorSqlScaningResultSelect
	ErrorDbNotCreaterd
	ErrorDbOpening
	ErrorDBCreateTable
	ErrorDBInsertPrepare
	ErrorDBRunningInsert
	ErrorK8sClusterNotFoundInKubeconfig
	ErrorK8sAuthInfoNotFoundInKubeconfig
	ErrorK8sLoadingKubeconfig
	ErrorGeneric
	ErrorK8sConfIsNotProperlyInitialized
	FailedToConnectToKubernetesAPIServer
)

func (k8s errorManager) BuildMsgError(params ...any) ErrorSystem {

	var errorSystem ErrorSystem

	switch k8s {

	case ErrorK8sObjectnotSupported:
		errorSystem = ErrorSystem{1, fmt.Sprintf("the k8s object %s is not supported", params[0])}
	case ErrorK8sGeneratingDynamicClient:
		errorSystem = ErrorSystem{2, fmt.Sprintf("creating dynamic client: %v", params[0])}
	case ErrorK8sRestConfig:
		errorSystem = ErrorSystem{3, fmt.Sprintf("generatig the rest conf for the path %s. %v", params[0], params[1])}
	case ErrorK8sClientConfig:
		errorSystem = ErrorSystem{4, fmt.Sprintf("generatig the client conf for the path %s. %v", params[0], params[1])}
	case ErrorK8sRestResource:
		errorSystem = ErrorSystem{5, fmt.Sprintf("getting the resource %v in the ns %s: %v", params[0], params[1], params[2])}
	case ErrorJsonMarshallResourceList:
		errorSystem = ErrorSystem{6, fmt.Sprintf("in the conversion of items returned of k8s to JSON for key %s", params[0])}
	case ErrorK8sKubeconfgNotAccesible:
		errorSystem = ErrorSystem{7, fmt.Sprintf("the kubeconfig is not accesible in this path: %s. %v", params[0], params[1])}
	case ErrorKubeAliasPathNotDefined:
		errorSystem = ErrorSystem{8, ("the env var KUBEALIAS is not defined")}
	case ErrorKubeAliasReadingFile:
		errorSystem = ErrorSystem{9, fmt.Sprintf("reading the kube_alias file in this path:%s. %v", params[0], params[1])}
	case ErrorKubeAliasVersionNotFoud:
		errorSystem = ErrorSystem{10, ("the version tag not found")}
	case ErrorKubeAliasNotFoud:
		errorSystem = ErrorSystem{11, fmt.Sprintf("the alias %s is not defined within the alias file", params[0])}
	case ErrorSqlRuningSelect:
		errorSystem = ErrorSystem{12, fmt.Sprintf(`failed executing SQL:" %s". Details:%v`, params[0], params[1])}
	case ErrorSqlReadingColumns:
		errorSystem = ErrorSystem{13, fmt.Sprintf("reading the columns:%v", params[0])}
	case ErrorSqlScaningResultSelect:
		errorSystem = ErrorSystem{14, fmt.Sprintf("failed to scan row: %v", params[0])}
	case ErrorKubeAliasParseFile:
		errorSystem = ErrorSystem{15, fmt.Sprintf("parsing the kube_alias file in this path:%s. %v", params[0], params[1])}
	case ErrorK8sRestResourceWithoutNS:
		errorSystem = ErrorSystem{16, fmt.Sprintf("getting the resource %v: %v", params[0], params[1])}
	case ErrorSqlNotASelect:
		errorSystem = ErrorSystem{17, fmt.Sprintf("the query is not a select: %s", params[0])}
	case ErrorDbNotCreaterd:
		errorSystem = ErrorSystem{18, fmt.Sprintf("creating the DB object: %v", params[0])}
	case ErrorDbOpening:
		errorSystem = ErrorSystem{18, fmt.Sprintf("opening the DB object: %v", params[0])}
	case ErrorDBCreateTable:
		errorSystem = ErrorSystem{19, fmt.Sprintf("creatiung the DB %s: %v", params[0], params[1])}
	case ErrorDBInsertPrepare:
		errorSystem = ErrorSystem{20, fmt.Sprintf("preparing the insert within the table %s: %v", params[0], params[1])}
	case ErrorDBRunningInsert:
		errorSystem = ErrorSystem{21, fmt.Sprintf("executing the insert within the table %s: %v", params[0], params[1])}
	case ErrorK8sNotContext:
		errorSystem = ErrorSystem{22, "current Context: Not explicitly set in kubeconfig."}
	case ErrorK8sAuthInfoNotFoundInKubeconfig:
		errorSystem = ErrorSystem{21, fmt.Sprintf("executing the insert within the table %s: %v", params[0], params[1])}
	case ErrorK8sClusterNotFoundInKubeconfig:
		errorSystem = ErrorSystem{22, fmt.Sprintf("Cluster '%s' not found in kubeconfig", params[0] ) }
	case ErrorK8sLoadingKubeconfig:
		errorSystem = ErrorSystem{22, fmt.Sprintf("error loading kubeconfig from %s: %v", params[0], params[1] ) }
	case ErrorGeneric:
		errorSystem = ErrorSystem{22, fmt.Sprintf("generic error: %v", params[0] ) }
	case ErrorK8sConfIsNotProperlyInitialized:
		errorSystem = ErrorSystem{22, "K8sConf is not properly initialized: restConf or clientConf is nil" }
	
	case FailedToConnectToKubernetesAPIServer:
		errorSystem = ErrorSystem{22, fmt.Sprintf("failed to connect to Kubernetes API server: %v", params[0] ) }
	}

	errorSystem.errorMsg = fmt.Sprintf("error msg: %s\nerror code: %d", errorSystem.errorMsg, errorSystem.errorCode)
	return errorSystem
}

func (err ErrorSystem) KO() {
	utils.Logger(utils.ERROR, err.errorMsg)
	os.Exit(err.errorCode)
}
