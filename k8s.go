package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func (resource defaultResource) retrieveContent(restConf *rest.Config, key string, channel chan<- map[string]string) {
	dynamicClient, err := dynamic.NewForConfig(restConf)
	if err != nil {
		ErrorK8sGeneratingDynamicClient.buildMsgError(err).KO()
	}
	resourceList, err := dynamicClient.Resource(resource.GroupVersionResource).Namespace(resource.NameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil{
		ErrorK8sRestResource.buildMsgError(resource.NameSpace, resource.GroupVersionResource).KO()
	}
	if conentBytes, err := json.Marshal(resourceList.Items); err != nil {
		ErrorJsonMarshallResourceList.buildMsgError(key).KO()
	} else {
		channel <- map[string]string{
			key: string(conentBytes),
		}
	}
}

type K8sConf struct {
	restConf   *rest.Config
	clientConf *kubernetes.Clientset
}

/*
createConfiguration given a path, generate a K8sconf to store the client and rest configuration. If the path is empty, then will verify the env var KUBECONFIG,
if is also empty then it will check the default path.
*/
func createConfiguration(path string) K8sConf {
	if path != "" {
		os.Setenv("KUBECONFIG", path)
	}
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	}
	k8sConf := K8sConf{}
	if restConf, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath); err != nil {
		ErrorK8sRestConfig.buildMsgError(kubeconfigPath).KO()
	} else {
		k8sConf.restConf = restConf
	}
	if clientConfig, err := kubernetes.NewForConfig(k8sConf.restConf); err != nil {
		ErrorK8sClientConfig.buildMsgError(kubeconfigPath, err).KO()
	} else {
		k8sConf.clientConf = clientConfig
	}
	return k8sConf
}

// generateMapObjects retrieve from cluster the map of Resource by name and alias.
func generateMapObjects(clientConfig *kubernetes.Clientset, ns string) map[string]defaultResource {

	result := map[string]defaultResource{}

	//Retrieve the list of apiResources
	apiResourceLists, _ := clientConfig.Discovery().ServerPreferredResources()

	// Iterate over each ApiResource and their resource
	for _, apiResourceList := range apiResourceLists {
		for _, apiResource := range apiResourceList.APIResources {
			defaultResource := defaultResource{
				GroupVersionResource: schema.GroupVersionResource{
					Version:  apiResourceList.GroupVersion,
					Group:    apiResource.Group,
					Resource: apiResource.Name,
				}, NameSpace: ns,
			}
			if len(apiResource.ShortNames) > 0 {
				// Check if there some aliases, in that case, for each alias it will store a new entry
				for _, alias := range apiResource.ShortNames {
					result[alias] = defaultResource
				}
			}
			result[apiResource.SingularName] = defaultResource
			result[apiResource.Name] = defaultResource
		}
	}
	return result
}

type defaultResource struct {
	schema.GroupVersionResource
	NameSpace string
}

type kubeParams struct {
	namespace    string
	k8sObjs      []string
	kubeconfPath string
}

func RetrieveK8sObjects(params kubeParams) map[string]string {

	conf := createConfiguration(params.kubeconfPath)
	mapK8sObject := generateMapObjects(conf.clientConf, params.namespace)

	chanResult := make(chan map[string]string)
	var wg sync.WaitGroup
	for _, keyObject := range params.k8sObjs {
		obj, ok := mapK8sObject[keyObject]
		if !ok {
			ErrorK8sObjectnotSupported.buildMsgError(keyObject).KO()
		}
		wg.Add(1)
		go obj.retrieveContent(conf.restConf, keyObject, chanResult)
	}

	wg.Wait()
	close(chanResult)

	fmt.Println(mapK8sObject)

	result, ok := <-chanResult
	if ok {
		fmt.Print("Channep open")
	} else {
		fmt.Println(" Clannel closed")
	}

	return result
}
