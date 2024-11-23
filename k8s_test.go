package main

import "testing"

func Test_createConfiguration(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Add default kubeconf", args{path: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createConfiguration(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("createConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
