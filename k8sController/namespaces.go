package k8sController

import (
	"context"
	"errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	log "k8s.io/klog/v2"

	"runtime"
)

type NameSpace struct {
	Name string
	Pods []*v1.Pod
}

func GetPodsInfoByNS(nsName string) (*v1.PodList, error) {
	_, kc := NewClient()
	exist := NamespaceExist(nsName)
	if !exist {
		return nil, nil
	}
	podList, err := kc.ClientSet.CoreV1().Pods(nsName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podList, nil
}
func GetAllNamespace() (*v1.NamespaceList, error) {
	_, kc := NewClient()
	nsList, err := kc.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nsList, nil
}
func NamespaceExist(ns string) bool {

	nsList, err := GetAllNamespace()
	if err != nil {
		return false
	}

	for _, nowNS := range nsList.Items {
		if ns == nowNS.Name {
			return true
		}
	}
	return false

}
func CreateNamespace(ns string) error {
	if ns == "" {
		err := errors.New("ns is nil")
		log.Warning(err.Error())
		return err
	}
	if NamespaceExist(ns) {
		return nil
	}
	log.Infof("create " + ns + "starting")
	kc := K8sClient()
	opts := metav1.CreateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
		DryRun:       nil,
		FieldManager: "",
	}
	listOptions := metav1.ListOptions{}
	nss, err := kc.CoreV1().Namespaces().List(context.TODO(), listOptions)

	for _, oldNS := range nss.Items {
		if oldNS.Namespace == ns {
			log.Infof("[CreateNamespace] %s is exist", ns)
			return nil
		}
	}
	nsmo := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
	}

	_, err = kc.CoreV1().Namespaces().Create(context.TODO(), nsmo, opts)
	if err != nil {
		return err
	}

	log.Infof("[CreateNamespace] %s is Create", ns)
	return nil
}
func DeleteNamespaces(ns string) error {
	if !NamespaceExist(ns) {
		log.Infof("[DeleteNamespaces] %s is not exist", ns)

		return nil
	}
	defer runtime.GC()
	kc := K8sClient()
	var opts metav1.DeleteOptions
	if err := kc.CoreV1().Namespaces().Delete(context.TODO(), ns, opts); err != nil {
		log.Warning(err)
		return err
	}
	log.Infof("[DeleteNamespaces] %s is Delete", ns)

	return nil

}
