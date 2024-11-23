package main

import (
	"context"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)



func createConfiguration(path string) ( err error) {

	if home := homedir.HomeDir(); home != ""{
		path = filepath.Join(home, ".kube", "config")
	}



	config , err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil{
		return  fmt.Errorf("error getting kubeconfig conf %d" , err)
	}


	 clientConfig, _ := kubernetes.NewForConfig(config)
	 core := clientConfig.CoreV1()
	cms, _ := core.ConfigMaps("").List(context.TODO(),  metav1.ListOptions{})

	for _, item :=  range cms.Items{
		fmt.Printf("Name: %s, ns: %s\n", item.Name, item.Namespace)
	}

	 apiResources, _ := clientConfig.Discovery().ServerPreferredResources()

	 for _, lstApiResources := range apiResources{
		for _ , resource := range lstApiResources.APIResources{
			fmt.Printf("%s %s \n", resource.Name, resource.ShortNames)
		}
	 } 

	 return nil

}


