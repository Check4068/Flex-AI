package plugin

import (
	"fmt"

	"volcano.sh/volcano/pkg/scheduler/api"
	"volcano.sh/volcano/pkg/scheduler/plugins/xpu-scheduler-plugin/util"
)

func (sh *ScheduleHandler) initXPUDevicesOfNode(sJob *SchedulerJob, node *api.NodeInfo) {
	xpuDevices := sJob.handler.GetXPUDevicesFromNode(node)
	sh.Lock()
	sh.XPUDevices[node.Name] = xpuDevices
	sh.Unlock()
}

func (sh *ScheduleHandler) NodePredicate(task *api.TaskInfo, node *api.NodeInfo) ([]*api.Status, error) {
	if sh == nil || task == nil || node == nil {
		return nil, fmt.Errorf("invalid input")
	}
	predicateStatus := make([]*api.Status, 0)
	sJob, ok := sh.Jobs[task.Job]
	if !ok {
		return predicateStatus, nil
	}
	if !util.IsXPUName(sJob.ReqXPUName) || !IsXPUTask(sJob, task) {
		return predicateStatus, nil
	}
	if err := sJob.preCheckNodePredicate(task, node); err != nil {
		checkStatus := &api.Status{
			Code:   api.Unschedulable,
			Reason: err.Error(),
		}
		predicateStatus = append(predicateStatus, checkStatus)
		return predicateStatus, err
	}

	code, err := sJob.handler.NodePredicateForTask(sJob, task, node, sh)
	if err != nil {
		checkStatus := &api.Status{
			Code:   code,
			Reason: err.Error(),
		}
		predicateStatus = append(predicateStatus, checkStatus)
		return predicateStatus, err
	}
	return predicateStatus, nil
}
