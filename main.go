package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	collection "github.com/jose78/go-collections"
	"gopkg.in/yaml.v2"
)

const (
	CTX_KEY_ALIAS_NAME = "alias.name"
	CTX_KEY_ALIAS_ARGS = "alias.args"
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

func destroy(){
	sqliteDatabase.Close()
	os.Remove("./sqlite-database.db")
}

func main() {
	ctx := context.Background()


	ctx = context.WithValue(ctx, CTE_NS, "")

	if len(os.Args) > 2 {
		ctx = context.WithValue(ctx, CTX_KEY_ALIAS_ARGS, os.Args[2:])
	}
	ctx = context.WithValue(ctx, CTX_KEY_ALIAS_NAME, os.Args[1])

	loadKubeAlias().execute(ctx)
	defer destroy() 

	//k8sParams :=  kubeParams{namespace: "", kubeconfPath: "", k8sObjs: []string{} }
	//	objRetrieved := RetrieveK8sObjects(k8sParams)
	//
	//
	//
	//
	//
	//	objRetrieved, _ := k8sGetElements(nil)
	//	for table, value := range objRetrieved {
	//		createTable(sqliteDatabase, table) // Create Database Tables
	//		insert(sqliteDatabase, value, table)
	//	}

	// 	sqlSelect := `SELECT
	//     json_extract(p.data, '$.metadata.name') AS pod_name,
	//     json_extract(d.data, '$.metadata.name') AS deploy_name
	// FROM
	//     pod AS p
	// JOIN
	//     deploy AS d ON json_extract(p.data, '$.metadata.labels.app') = json_extract(d.data, '$.metadata.name')
	// ;`

}

type Command interface {
	execute(context.Context)
}

func factoryCommand(version string) func([]byte) Command {
	var fn func([]byte) Command
	version = strings.TrimSpace(version)
	switch version {
	case "v1":
		fn = func(s []byte) Command {
			alias := AliasDefV1{}
			yaml.Unmarshal(s, &alias)
			return alias
		}
	}

	return fn
}

// loadKubeAlias
func loadKubeAlias() Command {

	kubepath := os.Getenv("KUBEALIAS")
	if kubepath == "" {
		ErrorKubeAliasPathNotDefined.buildMsgError().KO()
	}

	aliasByteConent, errReadinfKubeakis := os.ReadFile(kubepath)
	if errReadinfKubeakis != nil {
		ErrorKubeAliasReadingFile.buildMsgError(kubepath, errReadinfKubeakis)
	}

	// Extract the version declared within the file
	constent := strings.Split(string(aliasByteConent), "\n")
	var result []string
	collection.Filter(func(data string) bool { return strings.HasPrefix(data, "version") }, constent, &result)

	// if version is not declared, then the application will be finished
	if len(result) == 0 {
		ErrorKubeAliasVersionNotFoud.buildMsgError().KO()
	}

	version := strings.Split(result[0], ":")[1]
	cmd := factoryCommand(version)(aliasByteConent)
	return cmd
}
