package main

import (
	"context"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	CTE_KUBECONFIG = "KUBECONFIG"
	CTE_NS         = "namespace"
	CTE_TABLE      = "tables"
)

func (resource defaultResource) retrieveContent(restConf *rest.Config, key string) []unstructured.Unstructured {
	dynamicClient, err := dynamic.NewForConfig(restConf)
	if err != nil {
		ErrorK8sGeneratingDynamicClient.buildMsgError(err).KO()
	}

	var resourceList *unstructured.UnstructuredList
	var errResource error
	if resource.NameSpace == "" {
		resourceList, errResource = dynamicClient.Resource(resource.GroupVersionResource).List(context.TODO(), metav1.ListOptions{})
	} else {
		resourceList, errResource = dynamicClient.Resource(resource.GroupVersionResource).Namespace(resource.NameSpace).List(context.TODO(), metav1.ListOptions{})
	}

	if errResource != nil {
		ErrorK8sRestResource.buildMsgError(resource.NameSpace, resource.GroupVersionResource).KO()
	}
	return resourceList.Items
}

type K8sConf struct {
	restConf   *rest.Config
	clientConf *kubernetes.Clientset
}

// retrieveKubeConf discover which is the path of kubeconfig
func retrieveKubeConf(ctx context.Context) string {
	path := ctx.Value(CTE_KUBECONFIG)
	if path != nil && path.(string) != "" {
		os.Setenv(CTE_KUBECONFIG, path.(string))
	}
	kubeconfigPath := os.Getenv(CTE_KUBECONFIG)
	if kubeconfigPath == "" {
		kubeconfigPath = clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	}
	return kubeconfigPath
}

/*
createConfiguration given a path, generate a K8sconf to store the client and rest configuration. If the path is empty, then will verify the env VAR_varKUBECONFIG
if is also empty then it will check the default path.
*/
func createConfiguration(path string) K8sConf {
	if path != "" {
		os.Setenv(CTE_KUBECONFIG, path)
	}
	kubeconfigPath := os.Getenv(CTE_KUBECONFIG)
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

// retrieveK8sObjects retrieve from k8s ckuster a map of list of componentes deployed
func retrieveK8sObjects(ctx context.Context) []unstructured.Unstructured {
	pathK8s := retrieveKubeConf(ctx)
	conf := createConfiguration(pathK8s)
	ns := ctx.Value(CTE_NS).(string)
	table := ctx.Value(CTE_TABLE).(string)
	mapK8sObject := generateMapObjects(conf.clientConf, ns)
	result := []unstructured.Unstructured{}

	obj, ok := mapK8sObject[table]
	if !ok {
		ErrorK8sObjectnotSupported.buildMsgError(table).KO()
	}
	func(conf K8sConf, table string) {
		k8sObjs := obj.retrieveContent(conf.restConf, table)
		result = k8sObjs
	}(conf, table)

	return result
}
