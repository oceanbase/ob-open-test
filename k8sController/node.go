package k8sController

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	log "k8s.io/klog/v2"
)

func NodeAddOBLabel() error {
	//get nodes
	kc := K8sClient()
	nodes, err := kc.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	//add label:    topology.kubernetes.io/zone: zone1
	for _, node := range nodes.Items {
		labels := node.Labels
		target, ok := labels["topology.kubernetes.io/zone"]
		if target == "zone1" && ok {
			continue
		}
		labels["topology.kubernetes.io/zone"] = "zone"
		node.SetLabels(labels)
	}
	return nil

}

func GetPodNameByNS(ns string) ([]string, error) {
	kc := K8sClient()
	var nodeNames []string
	opts := &metav1.ListOptions{}
	nodes, err := kc.CoreV1().Pods(ns).List(context.TODO(), *opts)
	if err != nil {
		log.Warning(err)
		return nil, err
	}
	for _, node := range nodes.Items {
		nodeNames = append(nodeNames, node.Name)
	}
	return nodeNames, nil
}
