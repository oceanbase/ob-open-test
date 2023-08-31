package k8sController

import (
	"context"
	"strings"
	"time"

	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	log "k8s.io/klog/v2"
)

func KillPod(ns, podName string) error {
	kc := K8sClient()
	opts := metav1.DeleteOptions{}
	err := kc.CoreV1().Pods(podName).Delete(context.TODO(), podName, opts)
	log.Warning(err)
	return err
}
func GetPodStatus(pod *coreV1.Pod) string {

	for _, cond := range pod.Status.Conditions {
		if string(cond.Type) == ContainersReady {
			if string(cond.Status) != ConditionTrue {
				return "Unavailable"
			}
		} else if string(cond.Type) == PodInitialized && string(cond.Status) != ConditionTrue {
			return "Initializing"
		} else if string(cond.Type) == PodReady {
			if string(cond.Status) != ConditionTrue {
				return "Unavailable"
			}
			for _, containerState := range pod.Status.ContainerStatuses {
				if !containerState.Ready {
					return "Unavailable"
				}
			}
		} else if string(cond.Type) == PodScheduled && string(cond.Status) != ConditionTrue {
			return "Scheduling"
		}
	}
	return string(pod.Status.Phase)
}

func WaitPodBeRunning(ns, podname string, duration time.Duration) bool {
	kc := K8sClient()
	opts := metav1.ListOptions{}
	timeout := time.Now().Unix()
	for {
		if (timeout - time.Now().Unix()) < 0 {
			return false
		}
		podList, err := kc.CoreV1().Pods(ns).List(context.TODO(), opts)
		if err != nil {
			return false
		}
		for _, pod := range podList.Items {
			if pod.Name == podname && pod.Status.Phase == coreV1.PodRunning {
				return true
			}
		}

		time.Sleep(50 * time.Millisecond)

	}
}

//PodExist podName Containsis ok
func PodExist(ns, podName string) bool {
	kc := K8sClient()
	opts := metav1.ListOptions{}
	podlist, err := kc.CoreV1().Pods(ns).List(context.TODO(), opts)
	if err != nil {
		log.Warning(err)
		return false
	}
	for _, pod := range podlist.Items {
		if strings.Contains(pod.Name, podName) {
			return true
		}

	}
	return false
}
