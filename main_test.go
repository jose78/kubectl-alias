package main

import (
	"os"
	"testing"
)

func Test_selectQuery(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"retrieve the query" , args{name: "ENV_EXIST"} , "data" , false },
		{"retrieve the error" , args{name: "NOT_ENV_EXIST"} , "" , true },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ! tt.wantErr{
				os.Setenv("K_FCK_" + tt.args.name, tt.want)
			}
			got, err := selectQuery(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("selectQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("selectQuery() = %v, want %v", got, tt.want)
			}
			if ! tt.wantErr{
				os.Unsetenv("K_FCK_" + tt.args.name)
			}
		})
	}
}
