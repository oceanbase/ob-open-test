package k8sController

const (
	ContainersReady string = "ContainersReady"
	PodInitialized  string = "Initialized"
	PodReady        string = "Ready"
	PodScheduled    string = "PodScheduled"
)

const (
	ConditionTrue    string = "True"
	ConditionFalse   string = "False"
	ConditionUnknown string = "Unknown"
)
