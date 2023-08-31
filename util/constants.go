package util

const (
	OBOperatorNameBase   = "oceanbase-system"
	GetOBServerTagURL    = "https://hub.docker.com/v2/repositories/oceanbasedev/oceanbase-cn/tags/?page_size=100"
	OBClusterServiceName = "svc-ob-test"
)

//yamlPath
const (
	OBcrdPath             = "ob-operator/crd.yaml"
	OBOperatorPath        = "ob-operator/operator.yaml"
	OBClusterPath         = "ob-operator/obcluster.yaml"
	OBProxyDeploymentPath = "ob-operator/obproxy-Deployment.yaml"
	OBProxyServicePath    = "ob-operator/obproxy-Service.yaml"
	Ready                 = 3
	Running               = 4
	Error                 = 5
	NotAllReady           = 6
	//string consts
	OBoperatorNamespacs = "oceanbase-system"
	OBoperatorPodName   = "ob-operator-controller-manager"
	OBCrdName           = "obclusters.cloud.oceanbase.com"
	OBClusterYamlPath   = "ob-operator/obcluster.yaml"
)

//OBClusterState
const (
	OBClusterStateSysErr  = "SysErr"
	OBClusterStateUnready = "Unready"
	OBClusterStateReady   = "Ready"
)
