package main

import (
	"fmt"
	"os"
	"strings"

	collection "github.com/jose78/go-collections"
)

func main() {

}

func mapper(item string) any {
	stringSplited := strings.Split(item, "=")
	value := ""
	key := stringSplited[0]
	if len(stringSplited) > 1 {
		value = stringSplited[1]
	}
	return collection.Touple{Key: key, Value: value}
}

// function to extract the query to be executed
func selectQuery(nameVar string ) (string, error) {
	envNameVar := fmt.Sprintf("K_FCK_%s"  ,  strings.ToUpper(nameVar))
	var predicate collection.Predicate[collection.Touple] = func(item collection.Touple) bool {
		result := strings.Compare(item.Key.(string), envNameVar) == 0 
		return result
	}

	lstKeysEnv := os.Environ()
	mapUpdated := map[string]string{}
	collection.Map(mapper, lstKeysEnv, mapUpdated)
	result := map[string]string{}
	collection.Filter(predicate, mapUpdated, result)

	if len(result) == 0 {
		return "", fmt.Errorf(`Error: not found env var %s`, envNameVar)
	} else {
		return result[envNameVar], nil
	}
}
