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
package database

import (
	"os"
	"testing"

	"github.com/jose78/kubectl-alias/commons"
	_ "github.com/mattn/go-sqlite3"
)

func Test_getDbFaile(t *testing.T) {
	tests := []struct {
		name string
		kubealias string
		want string
	}{
		{"Get path of db File" , "~/.kube/alias.yml" ,  "~/.kube/sqlite-database.db" }, 
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(commons.ENV_VAR_KUBEALIAS_NAME, tt.kubealias)
			if got := getDbFaile(); got != tt.want {
				t.Errorf("getDbFaile() = %v, want %v", got, tt.want)
			}
		})
	}
}
