package k8sController

import (
	"context"
	"io/ioutil"

	"github.com/pytimer/k8sutil/apply"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	log "k8s.io/klog/v2"
)

// K8sClient func will be abandoned
func K8sClient() *kubernetes.Clientset {
	_, kc := NewClient()
	return kc.ClientSet
}
func ApplyByYamlStr(cf *rest.Config, YamlStr string) error {
	dynamicClient, err := dynamic.NewForConfig(cf)
	if err != nil {
		log.Warning(err)
		return err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(cf)
	if err != nil {
		log.Warning(err)
		return err
	}
	log.Info("to apply")
	applyOptions := apply.NewApplyOptions(dynamicClient, discoveryClient)
	if err := applyOptions.Apply(context.TODO(), []byte(YamlStr)); err != nil {
		log.Warning(err)
		return err
	}
	return nil
}
func ApplyPath(cf *rest.Config, FilePath string) error {
	applyStr, err := ioutil.ReadFile(FilePath)
	if err != nil {
		log.Warning(err)
		return err
	}
	dynamicClient, err := dynamic.NewForConfig(cf)
	if err != nil {
		log.Warning(err)
		return err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(cf)
	if err != nil {
		log.Warning(err)
		return err
	}

	applyOptions := apply.NewApplyOptions(dynamicClient, discoveryClient)
	if err := applyOptions.Apply(context.TODO(), applyStr); err != nil {
		log.Warning(err)
		return err
	}
	return nil
}

func DeleteByFilePath(FilePath string) error {
	return nil
}
