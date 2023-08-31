package k8sController

import (
	"context"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	log "k8s.io/klog/v2"

	"os"
	"path/filepath"

	cloudv1 "github.com/oceanbase/ob-operator/apis/cloud/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	ClientSet       *kubernetes.Clientset
	DynamicClient   dynamic.Interface
	DiscoveryClient *discovery.DiscoveryClient
}

func NewClient() (*rest.Config, *Client) {
	client := new(Client)
	if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".kube", "config")); err != nil {
		config, _ := rest.InClusterConfig()
		client.ClientSet, _ = kubernetes.NewForConfig(config)
		client.DynamicClient, _ = dynamic.NewForConfig(config)
		client.DiscoveryClient, _ = discovery.NewDiscoveryClientForConfig(config)
		return config, client
	} else {
		filePath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, _ := clientcmd.BuildConfigFromFlags("", filePath)
		client.ClientSet, _ = kubernetes.NewForConfig(config)
		client.DynamicClient, _ = dynamic.NewForConfig(config)
		client.DiscoveryClient, _ = discovery.NewDiscoveryClientForConfig(config)
		return config, client
	}
}

func (client *Client) GetResource() ([]*metav1.APIGroup, []*metav1.APIResourceList) {
	group, source, _ := client.DiscoveryClient.ServerGroupsAndResources()
	return group, source
}

func (client *Client) GetKind(kind string) string {
	_, resourceList := client.GetResource()
	for _, list := range resourceList {
		for _, resource := range list.APIResources {
			if resource.Kind == kind {
				return resource.Name
			}
		}
	}
	return ""
}

func (client *Client) GetGVR(unStruct *unstructured.Unstructured) *schema.GroupVersionResource {
	gvk := unStruct.GroupVersionKind()
	kind := client.GetKind(gvk.Kind)
	return &schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: kind,
	}
}

func (client *Client) CreateObj(obj unstructured.Unstructured) error {
	gvr := client.GetGVR(&obj)
	_, err := client.DynamicClient.Resource(*gvr).Namespace(obj.GetNamespace()).Create(context.TODO(), &obj, metav1.CreateOptions{})
	if err != nil {
		log.Info(err)
		return err
	}
	return nil
}
func (client *Client) GetOBClusterStatus(namespace, name string) string {
	instance, err := client.GetOBClusterInstance(namespace, name)
	if err == nil {
		return instance.Status.Status
	}
	return ""
}

var (
	OBClusterRes = schema.GroupVersionResource{
		Group:    OBClusterGroup,
		Version:  OBClusterVersion,
		Resource: OBClusterResource,
	}
)

const (
	OBClusterGroup    = "cloud.oceanbase.com"
	OBClusterVersion  = "v1"
	OBClusterKind     = "OBCluster"
	OBClusterResource = "obclusters"
)

func (client *Client) GetOBClusterInstance(namespace, name string) (cloudv1.OBCluster, error) {
	var instance cloudv1.OBCluster
	obj, err := client.DynamicClient.Resource(OBClusterRes).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		log.Info(err)
		return instance, err
	}
	_ = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &instance)
	return instance, nil
}
func (client *Client) GetObj(obj unstructured.Unstructured) (interface{}, error) {
	gvr := client.GetGVR(&obj)
	res, err := client.DynamicClient.Resource(*gvr).Namespace(obj.GetNamespace()).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
	if err != nil {
		log.Info(err)
		return res, err
	}
	return res, nil
}

func (client *Client) UpdateObj(obj unstructured.Unstructured) error {
	gvr := client.GetGVR(&obj)
	_, err := client.DynamicClient.Resource(*gvr).Namespace(obj.GetNamespace()).Update(context.TODO(), &obj, metav1.UpdateOptions{})
	if err != nil {
		log.Info(err)
		return err
	}
	return nil
}

func (client *Client) DeleteObj(obj unstructured.Unstructured) {
	gvr := client.GetGVR(&obj)
	_ = client.DynamicClient.Resource(*gvr).Namespace(obj.GetNamespace()).Delete(context.TODO(), obj.GetName(), metav1.DeleteOptions{})
}

func GetObjFromYaml(filePath string) unstructured.Unstructured {
	obj := MakeObjectFromFile(filePath)
	return obj.(unstructured.Unstructured)
}
func MakeObjectFromFile(filePath string) interface{} {
	yamlFile, _ := ioutil.ReadFile(filePath)
	obj := &unstructured.Unstructured{}
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, _, _ = dec.Decode(yamlFile, nil, obj)
	return *obj
}
