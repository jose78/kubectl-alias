package service

import (
	"context"
	"os"
	"strings"

	collection "github.com/jose78/go-collections"
	"github.com/jose78/kubectl-fuck/commons"
	"github.com/jose78/kubectl-fuck/internal/alias"
	"gopkg.in/yaml.v2"
)

func runAlias() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, commons.CTE_NS, "")

	if len(os.Args) > 2 {
		ctx = context.WithValue(ctx, commons.CTX_KEY_ALIAS_ARGS, os.Args[2:])
	}
	ctx = context.WithValue(ctx, commons.CTX_KEY_ALIAS_NAME, os.Args[1])

	loadKubeAlias().Execute(ctx)
}

func factoryCommand(version string) func([]byte) alias.Command {
	var fn func([]byte) alias.Command
	version = strings.TrimSpace(version)
	switch version {
	case "v1":
		fn = func(s []byte) alias.Command {
			alias := alias.AliasDefV1{}
			yaml.Unmarshal(s, &alias)
			return alias
		}
	}
	return fn
}

// loadKubeAlias
func loadKubeAlias() alias.Command {

	kubepath := os.Getenv("KUBEALIAS")
	if kubepath == "" {
		commons.ErrorKubeAliasPathNotDefined.BuildMsgError().KO()
	}

	aliasByteConent, errReadinfKubeakis := os.ReadFile(kubepath)
	if errReadinfKubeakis != nil {
		commons.ErrorKubeAliasReadingFile.BuildMsgError(kubepath, errReadinfKubeakis)
	}

	// Extract the version declared within the file
	constent := strings.Split(string(aliasByteConent), "\n")
	var result []string
	collection.Filter(func(data string) bool { return strings.HasPrefix(data, "version") }, constent, &result)

	// if version is not declared, then the application will be finished
	if len(result) == 0 {
		commons.ErrorKubeAliasVersionNotFoud.BuildMsgError().KO()
	}

	version := strings.Split(result[0], ":")[1]
	cmd := factoryCommand(version)(aliasByteConent)
	return cmd
}
