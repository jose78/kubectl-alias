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

package k8s

import (
	"context"
	"fmt"
	"os"

	"github.com/jose78/kubectl-alias/commons"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func (resource defaultResource) retrieveContent(restConf *rest.Config) []unstructured.Unstructured {
	dynamicClient, err := dynamic.NewForConfig(restConf)
	if err != nil {
		commons.ErrorK8sGeneratingDynamicClient.BuildMsgError(err).KO()
	}

	var resourceList *unstructured.UnstructuredList
	var errResource error
	if resource.NameSpace == "" {
		resourceList, errResource = dynamicClient.Resource(resource.GroupVersionResource).List(context.TODO(), metav1.ListOptions{})
		if errResource != nil {
			commons.ErrorK8sRestResourceWithoutNS.BuildMsgError(resource, errResource).KO()
		}
	} else {
		resourceList, errResource = dynamicClient.Resource(resource.GroupVersionResource).Namespace(resource.NameSpace).List(context.TODO(), metav1.ListOptions{})
		if errResource != nil {
			commons.ErrorK8sRestResource.BuildMsgError(resource, resource.NameSpace, errResource).KO()
		}
	}

	return resourceList.Items
}

type K8sConf struct {
	restConf   *rest.Config
	clientConf *kubernetes.Clientset
}

// retrieveKubeConf discover which is the path of kubeconfig
func retrieveKubeConf(path string) string {
	if path != "" {
		os.Setenv(commons.ENV_VAR_KUBECONFIG, path)
	}
	kubeconfigPath := os.Getenv(commons.ENV_VAR_KUBECONFIG)
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
		os.Setenv(commons.ENV_VAR_KUBECONFIG, path)
	}
	kubeconfigPath := os.Getenv(commons.ENV_VAR_KUBECONFIG)
	if kubeconfigPath == "" {
		kubeconfigPath = clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	}
	k8sConf := K8sConf{}
	if restConf, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath); err != nil {
		commons.ErrorK8sRestConfig.BuildMsgError(kubeconfigPath).KO()
	} else {
		k8sConf.restConf = restConf
	}
	if clientConfig, err := kubernetes.NewForConfig(k8sConf.restConf); err != nil {
		commons.ErrorK8sClientConfig.BuildMsgError(kubeconfigPath, err).KO()
	} else {
		k8sConf.clientConf = clientConfig
	}
	return k8sConf
}

type k8sConfig struct {
	pathK8sConfig string
	namespaceDefault string
	k8sResources map[string]defaultResource
}

// GenerateMapObjects retrieve from cluster the map of Resource by name and alias.
func GenerateMapObjects(config k8sConfig ) map[string]defaultResource {
	ns := config.namespaceDefault
	pathK8s := retrieveKubeConf(config.pathK8sConfig)
	conf := createConfiguration(pathK8s)
	clientConfig := conf.clientConf
	
	//Retrieve the list of apiResources
	apiResourceLists, _ := clientConfig.Discovery().ServerPreferredResources()
	result := map[string]defaultResource{}

	// Procesar los recursos
	for _, apiResourceList := range apiResourceLists {
		groupVersion, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if err != nil {
			fmt.Printf("Error parsing GroupVersion: %v", err)
			continue
		}

		// Iterar sobre los recursos individuales
		for _, resource := range apiResourceList.APIResources {
			defaultNs := ""
			if resource.Namespaced{
				defaultNs = ns
			}
			defaultResource := defaultResource{
					GroupVersionResource: schema.GroupVersionResource{
						Version:  groupVersion.Version,
						Group:    groupVersion.Group,
						Resource: resource.Name,
					}, NameSpace: defaultNs,
				}
			if len(resource.ShortNames) > 0 {
				// Check if there some aliases, in that case, for each alias it will store a new entry
				for _, alias := range resource.ShortNames {
					result[alias] = defaultResource
				}
			}
			result[resource.SingularName] = defaultResource
			result[resource.Name] = defaultResource
		}
	}
	return result
}

type defaultResource struct {
	schema.GroupVersionResource
	NameSpace string
}

// RetrieveK8sObjects retrieve from k8s ckuster a map of list of componentes deployed
func RetrieveK8sObjects(config k8sConfig , table string) []unstructured.Unstructured {
	
	pathK8s := retrieveKubeConf(config.pathK8sConfig)
	conf := createConfiguration(pathK8s)
	mapK8sObject :=  config.k8sResources
	result := []unstructured.Unstructured{}

	obj, ok := mapK8sObject[table]
	if !ok {
		commons.ErrorK8sObjectnotSupported.BuildMsgError(table).KO()
	}
	func(conf K8sConf, table string) {
		k8sObjs := obj.retrieveContent(conf.restConf)
		result = k8sObjs
	}(conf, table)

	return result
}
