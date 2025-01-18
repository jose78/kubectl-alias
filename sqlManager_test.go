package main

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_manipulateSElect(t *testing.T) {

	real_query := "SELECT ns.metadata.perconte.name AS ns_name,  p.metadata.name AS pod_name FROM ns, pod AS p WHERE  ns.metadata.namespace = p.metadata.name;"
	query_join := "SELECT p.id AS ID, p.name AS user_name FROM pepe p, fulgencio f WHERE p.name = 'tete' and p.id = f.id"
	query_simple := "SELECT p.id AS ID, p.name AS user_name FROM pepe p WHERE p.name = 'tete'"
	query_order_by := "SELECT p.id AS ID, p.name AS user_name FROM pepe p WHERE p.name = 'tete' ORDER BY p.name DESC"
	query_group_by := "SELECT p.id AS ID, p.name AS user_name FROM pepe p WHERE p.name = 'tete' GROUP BY p.type"
	query_group_by_and_order_by := "SELECT p.id AS ID, p.name AS user_name FROM pepe p WHERE p.name = 'tete' GROUP BY p.type ORDER BY p.name DESC"

	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Update a Query with Join", args{query_join}, "select json_extract(pepe, `\"$.Object.id\"`) as ID, json_extract(pepe, `\"$.Object.name\"`) as user_name from pepe as p, fulgencio as f where json_extract(pepe, `\"$.Object.name\"`) = 'tete' and json_extract(pepe, `\"$.Object.id\"`) = json_extract(fulgencio, `\"$.Object.id\"`)"},
		{"Update a simple Query", args{query_simple}, "select json_extract(pepe, `\"$.Object.id\"`) as ID, json_extract(pepe, `\"$.Object.name\"`) as user_name from pepe as p where json_extract(pepe, `\"$.Object.name\"`) = 'tete'"},
		{"Update a Query with Join", args{real_query}, "select json_extract(ns, `\"$.Object.metadata.perconte.name\"`) as ns_name, json_extract(pod, `\"$.Object.metadata.name\"`) as pod_name from ns, pod as p where json_extract(ns, `\"$.Object.metadata.namespace\"`) = json_extract(pod, `\"$.Object.metadata.name\"`)"},
		{"Update a Query with order by", args{query_order_by}, "select json_extract(pepe, `\"$.Object.id\"`) as ID, json_extract(pepe, `\"$.Object.name\"`) as user_name from pepe as p where json_extract(pepe, `\"$.Object.name\"`) = 'tete' order by json_extract(pepe, `\"$.Object.name\"`) desc"},
		{"Update a Query with group by", args{query_group_by}, "select json_extract(pepe, `\"$.Object.id\"`) as ID, json_extract(pepe, `\"$.Object.name\"`) as user_name from pepe as p where json_extract(pepe, `\"$.Object.name\"`) = 'tete' group by json_extract(pepe, `\"$.Object.type\"`)"},
		{"Update a complex Query", args{query_group_by_and_order_by}, "select json_extract(pepe, `\"$.Object.id\"`) as ID, json_extract(pepe, `\"$.Object.name\"`) as user_name from pepe as p where json_extract(pepe, `\"$.Object.name\"`) = 'tete' group by json_extract(pepe, `\"$.Object.type\"`) order by json_extract(pepe, `\"$.Object.name\"`) desc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapTables := findTablesWithAliases(tt.args.query)
			got := updateQuery(tt.args.query, mapTables)
			if got != tt.want {
				t.Errorf("manipulateSElect() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_regenerateColInfo(t *testing.T) {
//	type args struct {
//		col                  string
//		flagWithOutQualifier bool
//	}
//	tests := []struct {
//		name string
//		args args
//		want colInfo
//	}{
//		{"regenerate colInfo", args{col: "p____metadata____name", flagWithOutQualifier:  false},colInfo{name: "metadata.name", qualifier: "p"}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := regenerateColInfo(tt.args.col, tt.args.flagWithOutQualifier); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("regenerateColInfo() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
