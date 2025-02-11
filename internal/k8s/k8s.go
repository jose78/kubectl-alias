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

// retrieveContent retrieves the list of Kubernetes resources of the specified type,
// using the provided REST configuration to communicate with the Kubernetes API server.
//
// Parameters:
//   - restConf (*rest.Config): The Kubernetes REST configuration used to authenticate
//     and connect to the API server.
//
// Returns:
//   - []unstructured.Unstructured: A list of Kubernetes resources represented as
//     unstructured data. This format is useful for handling dynamic or unknown resource schemas.
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

// K8sConf holds the configuration and client needed to interact with a Kubernetes cluster.
type K8sConf struct {
	// restConf contains the Kubernetes REST configuration used to establish communication
	// with the Kubernetes API server. This includes authentication details, API server URL,
	// and other connection settings.
	restConf *rest.Config

	// clientConf is a Kubernetes clientset, which provides strongly-typed clients
	// for interacting with various Kubernetes resources (e.g., Pods, Deployments, Services).
	// It is built using the provided restConf.
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

// createConfiguration creates a Kubernetes configuration and clientset based on the provided kubeconfig file path.
//
// Parameters:
//   - pathKubeConfig (string): The file path to the kubeconfig file, which contains
//     the cluster, user, and authentication details needed to connect to the Kubernetes API server.
//
// Returns:
//   - K8sConf: A struct containing the REST configuration (`restConf`) and the Kubernetes clientset (`clientConf`),
//     allowing interaction with the Kubernetes cluster.
func createConfiguration(pathKubeCondif string) K8sConf {
	if pathKubeCondif != "" {
		os.Setenv(commons.ENV_VAR_KUBECONFIG, pathKubeCondif)
	}
	kubeconfigPath := os.Getenv(commons.ENV_VAR_KUBECONFIG)
	if kubeconfigPath == "" {
		kubeconfigPath = clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	}
	k8sConf := K8sConf{}
	if restConf, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath); err != nil {
		commons.ErrorK8sRestConfig.BuildMsgError(kubeconfigPath, err).KO()
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


// K8sInfo encapsulates information and configuration needed to interact with a Kubernetes cluster.
type K8sInfo struct {
	// PathK8sConfig specifies the path to the Kubernetes configuration file (kubeconfig).
	// This file is typically used to authenticate and connect to the Kubernetes API server.
	PathK8sConfig string

	// NamespaceDefault defines the default Kubernetes namespace to use when none is explicitly provided.
	// This ensures operations are scoped to the correct namespace by default.
	NamespaceDefault string

	// K8sResources is a map where the keys represent resource types (e.g., "pods", "services"),
	// and the values are defaultResource objects that define how to interact with these resources.
	// This map allows dynamic handling of Kubernetes resources based on their type.
	K8sResources map[string]defaultResource
}

// GenerateMapObjects generates a map of Kubernetes resource types to their corresponding defaultResource objects,
// using the provided Kubernetes configuration.
//
// Parameters:
//   - config (K8sInfo): A struct containing the necessary Kubernetes configuration, including
//     the path to the kubeconfig file, the default namespace, and existing resource mappings.
//
// Returns:
//   - map[string]defaultResource: A map where the keys represent resource types (e.g., "pods", "services"),
//     and the values are `defaultResource` objects that define how to interact with these resource types.
func GenerateMapObjects(info K8sInfo ) map[string]defaultResource {
	ns := info.NamespaceDefault
	pathK8s := retrieveKubeConf(info.PathK8sConfig)
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

// defaultResource represents a Kubernetes resource with its associated metadata and namespace information.
//
// Fields:
//   - GroupVersionResource (schema.GroupVersionResource): Specifies the group, version,
//     and resource type of the Kubernetes resource (e.g., apps/v1/replicasets).
//   - NameSpace (string): The namespace in which the resource resides. If empty, it typically
//     implies the resource is either cluster-wide or the default namespace is used.
type defaultResource struct {
	schema.GroupVersionResource
	NameSpace string
}

// RetrieveK8sObjects retrieves a list of Kubernetes resources of the specified type from the cluster,
// using the provided configuration.
//
// Parameters:
//   - config (K8sInfo): The Kubernetes configuration, including the path to the kubeconfig file,
//     the default namespace, and mappings of resource types.
//   - k8sObject (string): The name of the Kubernetes resource type to retrieve (e.g., "pods", "services").
//
// Returns:
//   - []unstructured.Unstructured: A list of unstructured Kubernetes resources of the specified type,
//     allowing for dynamic handling of their data structure.
//
// Notes:
//   - This function uses the information in `config` to determine the correct resource type and
//     namespace. If the specified resource type is not supported, an error or empty result may be returned.
func RetrieveK8sObjects(config K8sInfo , k8sObject string) []unstructured.Unstructured {
	
	pathK8s := retrieveKubeConf(config.PathK8sConfig)
	conf := createConfiguration(pathK8s)
	mapK8sObject :=  config.K8sResources
	result := []unstructured.Unstructured{}

	obj, ok := mapK8sObject[k8sObject]
	if !ok {
		commons.ErrorK8sObjectnotSupported.BuildMsgError(k8sObject).KO()
	}
	func(conf K8sConf) {
		k8sObjs := obj.retrieveContent(conf.restConf)
		result = k8sObjs
	}(conf)

	return result
}
