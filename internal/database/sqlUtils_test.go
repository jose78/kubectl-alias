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
	"reflect"
	"testing"
)

func TestWrapColumnsWithHisn(t *testing.T) {
	simple_select := "SELECT  ns.metadata.perconte.name AS ns_name,  p.metadata.name AS pod_name FROM ns, pod AS p WHERE  ns.metadata.namespace LIKE   'MEGANITO' ||  p.metadata.name;"
	simple_select_result := "select json_extract(ns , '$.Object.metadata.perconte.name') as ns_name, json_extract(pod , '$.Object.metadata.name') as pod_name from ns, pod as p where json_extract(ns , '$.Object.metadata.namespace') like 'MEGANITO'  ||  json_extract(pod , '$.Object.metadata.name')"
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Simple query with concatenation", args{query: simple_select}, simple_select_result, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ManipulateAST(tt.args.query, FindTablesWithAliases(tt.args.query))

			if got != tt.want {
				t.Errorf("ManipulateAST() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindTablesWithAliases(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"Find tables of a simple select without aliases", args{"select * from user, company"}, map[string]string{"user": "user", "company": "company"}},
		{"Find tables of a simple select witho aliases", args{"select * from user u, company c"}, map[string]string{"u": "user", "c": "company"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindTablesWithAliases(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindTablesWithAliases() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_regenerateColInfo(t *testing.T) {
	type args struct {
		colSplited   []string
		aliasToTable map[string]string
	}
	tests := []struct {
		name string
		args args
		want colInfo
	}{
		{"add prefix from name table" , args{colSplited: []string{ "events" , "regarding", "kind" },aliasToTable:  map[string]string{"events": "events", "p": "pods"}}, colInfo{columnName: "$.Object.regarding.kind", tableName: "events"} },
		{"add prefix without alias" , args{colSplited: []string{ "regarding", "kind" },aliasToTable:  map[string]string{"e": "events"}}, colInfo{columnName: "$.Object.regarding.kind", tableName: "events"} },
		{"add prefix from alias" , args{colSplited: []string{ "e" , "regarding", "kind" },aliasToTable:  map[string]string{"e": "events"}}, colInfo{columnName: "$.Object.regarding.kind", tableName: "events"} },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := regenerateColInfo(tt.args.colSplited, tt.args.aliasToTable); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("regenerateColInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
