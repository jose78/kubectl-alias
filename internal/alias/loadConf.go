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
	"os"
	"strings"

	"github.com/jose78/kubectl-alias/commons"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type builderCommand func(map[string]any) Alias

func FactoryAlias(aliasContent map[string]any) Alias {

	version, okVersionIsContained := aliasContent["version"]
	if !okVersionIsContained {
		commons.ErrorKubeAliasVersionNotFoud.BuildMsgError().KO()
	}

	var fn builderCommand
	version = strings.TrimSpace(version.(string))
	switch version {
	case "v1":
		fn = func(m map[string]any) Alias {
			var alias AliasDefV1
			mapstructure.Decode(m, &alias)
			return alias
		}
	}
	return fn(aliasContent)
}

func LoadKubeAlias() map[string]any {

	kubepath := os.Getenv(commons.ENV_VAR_KUBEALIAS_NAME)
	if kubepath == "" {
		commons.ErrorKubeAliasPathNotDefined.BuildMsgError().KO()
	}

	aliasByteConent, errReadingKubeAliasContent := os.ReadFile(kubepath)
	if errReadingKubeAliasContent != nil {
		commons.ErrorKubeAliasReadingFile.BuildMsgError(kubepath, errReadingKubeAliasContent).KO()
	}

	var result map[string]any
	errParsingKubeAliasContent := yaml.Unmarshal(aliasByteConent, &result)
	if errParsingKubeAliasContent != nil {
		commons.ErrorKubeAliasParseFile.BuildMsgError(kubepath, errParsingKubeAliasContent)
	}
	return result
}
