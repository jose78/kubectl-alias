package main

import (
	"reflect"
	"testing"
)

func Test_createConfiguration(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want K8sConf
	}{
		{"Retrieve Conf", args{path: ""}, K8sConf{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createConfiguration(tt.args.path)
			if got.clientConf == nil || got.restConf == nil {
				t.Errorf("createConfiguration() = %v", got)
			}
		})
	}
}

func TestRetrieveK8sObjects(t *testing.T) {
	type args struct {
		params kubeParams
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"", args{params: kubeParams{namespace: "", k8sObjs: []string{"pod"}}} , map[string]string{} },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RetrieveK8sObjects(tt.args.params); 
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveK8sObjects() = %v, want %v", got, tt.want)
			}
		})
	}
}
