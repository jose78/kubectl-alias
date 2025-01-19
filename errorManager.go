package main

import (
	"fmt"
	"os"
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
	ErrorK8sClientConfig
	ErrorK8sKubeconfgNotAccesible
	ErrorJsonMarshallResourceList
	ErrorKubeAliasPathNotDefined
	ErrorKubeAliasReadingFile
	ErrorKubeAliasVersionNotFoud
	ErrorKubeAliasNotFoud
	ErrorKubeAliasDuplicated
	ErrorSqlRuningSelect
	ErrorSqlReadingColumns
	ErrorSqlScaningResultSelect
)

func (k8s errorManager) buildMsgError(params ...any) ErrorSystem {

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
		errorSystem = ErrorSystem{5, fmt.Sprintf("getting the namespace %s;resource %v", params[1], params[0])}
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
	}

	errorSystem.errorMsg = fmt.Sprintf("error msg: %s\nerror code: %d", errorSystem.errorMsg, errorSystem.errorCode)
	return errorSystem
}

func (err ErrorSystem) KO() {
	fmt.Print(err.errorMsg)
	os.Exit(err.errorCode)
}
