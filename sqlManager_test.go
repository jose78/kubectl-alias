package main

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_manipulateSElect(t *testing.T) {

	query_join :=   "SELECT p.id AS ID, p.name AS user_name FROM pepe p, fulgencio f WHERE p.name = 'tete' and p.id = f.id"
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
		{"Update a Query with order by", args{query_order_by}, "select json_extract(pepe, `\"$.Object.id\"`) as ID, json_extract(pepe, `\"$.Object.name\"`) as user_name from pepe as p where json_extract(pepe, `\"$.Object.name\"`) = 'tete' order by json_extract(pepe, `\"$.Object.name\"`) desc"},
		{"Update a Query with group by", args{query_group_by}, "select json_extract(pepe, `\"$.Object.id\"`) as ID, json_extract(pepe, `\"$.Object.name\"`) as user_name from pepe as p where json_extract(pepe, `\"$.Object.name\"`) = 'tete' group by json_extract(pepe, `\"$.Object.type\"`)"},
		{"Update a complex Query", args{query_group_by_and_order_by}, "select json_extract(pepe, `\"$.Object.id\"`) as ID, json_extract(pepe, `\"$.Object.name\"`) as user_name from pepe as p where json_extract(pepe, `\"$.Object.name\"`) = 'tete' group by json_extract(pepe, `\"$.Object.type\"`) order by json_extract(pepe, `\"$.Object.name\"`) desc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapTables := findTablesWithAliases(tt.args.query)
			if got := updateAST(tt.args.query, mapTables); got != tt.want {
				t.Errorf("manipulateSElect() = %v, want %v", got, tt.want)
			}
		})
	}
}
