package main

//import (
//	"os"
//	"reflect"
//	"testing"
//
//	_ "github.com/mattn/go-sqlite3"
//	"gopkg.in/yaml.v2"
//)

//func Test_selectQuery(t *testing.T) {
//	type args struct {
//		name string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    string
//		wantErr bool
//	}{
//		{"retrieve the query", args{name: "ENV_EXIST"}, "data", false},
//		{"retrieve the error", args{name: "NOT_ENV_EXIST"}, "", true},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if !tt.wantErr {
//				os.Setenv("K_FCK_"+tt.args.name, tt.want)
//			}
//			got, err := selectQuery(tt.args.name)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("selectQuery() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("selectQuery() = %v, want %v", got, tt.want)
//			}
//			if !tt.wantErr {
//				os.Unsetenv("K_FCK_" + tt.args.name)
//			}
//		})
//	}
//}
//
//func Test_evaluateQuery(t *testing.T) {
//	type args struct {
//		query string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    sqlContainer
//		wantErr bool
//	}{
//		{name: "Evaluate query", args: args{query: "SELECT id, name as name_user FROM users WHERE id = 10"},
//			want: sqlContainer{sqlFrom: []string{"users"}, sqlWhere: "id = 10", sqlSelect: map[string]string{"id": "id", "NAME": "NAME_USER"}}, wantErr: false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := evaluateQuery(tt.args.query)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("evaluateQuery() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("evaluateQuery() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_loadKubeAlias(t *testing.T) {
//
//	aliasContent, _ := os.ReadFile("resources/kube_alias.yml")
//	alias := AliasV1{}
//	yaml.Unmarshal(aliasContent, &alias)
//
//	type args struct {
//		configureTest func()
//	}
//	tests := []struct {
//		name string
//		args args
//		want Command
//	}{
//		{"Load alias v1 default kube_alias", args{ configureTest: func ()  {
//			os.Setenv("KUBEALIAS", "/home/jose/ws/kubectl-fuck/resources/kube_alias.yml")
//		}}, alias},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.args.configureTest()
//			got := loadKubeAlias()
//			if reflect.DeepEqual(got, tt.want) {
//				t.Errorf("loadKubeAlias() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//