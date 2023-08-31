//一个关于测试DB稳定性的子工具，流程可进行自定义
/*
稳定性压测仅支持k8s测试
*/

package main

import (
	obopentestOB "ob-open-test/obopentest-ob"
	obopentest_task "ob-open-test/obopentest-task"
	"time"

	log "k8s.io/klog/v2"
)

var CaseName = "stability"

//
func main() {
	t := obopentest_task.NewTask()
	t.SetName(CaseName)

	//create obcluster
	obstep := obopentest_task.NewStep()
	obstep.Name = "obstart"
	obstep.Type = obopentest_task.StepTypeSerial
	obm := obopentestOB.OBModel{Name: CaseName}
	obConfMap := make(map[string]string)
	obConfMap["OBClusterName"] = CaseName
	obConfMap["order"] = "create"
	obm.SetConf(obConfMap)
	obstep.AddModel(&obm)
	t.AddStep(*obstep)
	t.Start()
	log.Info("add obstep")

	//CreateMinTenant 2c4g
	ob := obopentestOB.NewOBCluster(CaseName)
	ob.Ping()
	ob.CreateMinTenant()

	log.Info("add obstep end")
	time.Sleep(1 * time.Minute)

	//add chaos
	chaosStep := obopentest_task.NewStep()
	chaosStep.Name = "chaosStep"
	chaosStep.Type = obopentest_task.StepTypeSerial

	cm := obopentest_task.NewMode("chaosblade")
	cmconf := make(map[string]string)
	cmconf["type"] = "killPod"
	cmconf["OBClusterName"] = CaseName
	cm.SetConf(cmconf)
	chaosStep.AddModel(cm)

	t.AddStep(*chaosStep)

	t.Start()
	t.Destory()
	return

}
