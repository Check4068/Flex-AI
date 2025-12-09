package plugin

import (
	"sync"

	""
	"k8s.io/apimachinery/pkg/types"
)

const PluginName = "huaweiXPU"

const (
	scoreWeight              = 100
	defaultSchedulingTaskNum = 1
	scoreSplitMinSize        = 5
)

type ScheduleHandler struct {
	XPUPlugins     map[string]XPUBuilder
	XPUDevices     map[string]map[int]*common.XPUDevice
	Jobs           map[api.JobID]*SchedulerJob
	DeleteJobInfos map[api.JobID]*JobInfo
	SessionID      types.UID
	Nodes          []*api.NodeInfo
	*sync.Mutex
}

type ContainerDevices []common.ContainerDevice

type PodDevices []ContainerDevices

type SchedulerJob struct {
	Id            api.JobID
	ReferenceName string
	NameSpace     string
	Annotation    map[string]string
	Selector      map[string]string
	Label         map[string]string
	UnschedulableReason
	handler     XPUSchedulerPlugin
	JobReadyTag bool
	*util.XPUJob
	TopologyAllocateOnce sync.Once
}

type UnschedulableReason struct {
	Reason map[string]string
	*sync.Mutex
}
