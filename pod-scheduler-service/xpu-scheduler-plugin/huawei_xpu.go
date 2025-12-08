/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package main

import (
	"errors"
	"strings"
	"sync"

	"k8s.io/klog/v2"
	"volcano.sh/volcano/pkg/scheduler/api"
	"volcano.sh/volcano/pkg/scheduler/framwork"
	"volcano.sh/volcano/pkg/scheduler/plugins/xpu-scheduler-plugin/common"
	"volcano.sh/volcano/pkg/scheduler/plugins/xpu-scheduler-plugin/internal/xpu"
	"volcano.sh/volcano/pkg/scheduler/plugins/xpu-scheduler-plugin/plugin"
	"volcano.sh/volcano/pkg/scheduler/plugins/xpu-scheduler-plugin/util"
)

const (
	// PluginName huawei xpu plugin name, it must be the same as the so package name
	PluginName = "huawei-xpu"
	// TopologyEnable topology setting
	TopologyEnable = "TopologyEnable"
	// NumaEnable numa setting
	NumaEnable = "NumaEnable"
	// TestEnable test mode setting
	TestEnable = "TestEnable"
	// XPUTopologyNodeList node list setting
	XPUTopologyNodeList = "XPUTopologyNodeList"
	// XPUTopologyNodeBandwidth bandwidth setting between nodes
	XPUTopologyNodeBandwidth = "XPUTopologyNodeBandwidth"
)

var (
	scheduleHandler *plugin.scheduleHandler
	once            sync.Once
)

// GetScheduleHandler implement the singleton pattern of the ScheduleHandler
func GetScheduleHandler() *plugin.ScheduleHandler {
	once.Do(func() {
		scheduleHandler = HandlerCreate()
	})
	return scheduleHandler
}

type huaweiXPUPlugin struct {
	// Scheduler for plugin args and its handler.
	Scheduler *plugin.ScheduleHandler
	// Arguments given for the plugin
	Arguments framework.Arguments
}

func (xp *huaweiXPUPlugin) Name() string {
	return PluginName
}

// New huawei xpu framework plugin
func New(arguments framework.Arguments) framework.Plugin {
	return &huaweiXPUPlugin{Scheduler: GetScheduleHandler(), Arguments: arguments}
}

func getCommonConfig(args framework.Arguments) {
	args.GetBool(&xpu.Config.TopologyEnable, TopologyEnable)
	args.GetBool(&xpu.Config.NumaEnable, NumaEnable)
	args.GetBool(&xpu.Config.TestEnable, TestEnable)
}

func getNodeBandwidthConf(args framework.Arguments) {
	argv, ok := args[XPUTopologyNodeList]
	if !ok {
		return
	}
	value, ok := argv.(string)
	if !ok {
		klog.V(util.LogErrorLevel).Infof("XPUTopologyNodeList in args is not string")
		return
	}
	tmp := strings.Split(value, util.Comma)
	topologyNodeList := tmp

	err := getNodeBandwidth(args, topologyNodeList)
	if err != nil {
		klog.V(util.LogErrorLevel).Infof("get node bandwidth failed, err: %v", err.Error())
		util.XPUTopologyNodeBandwidth = nil
	}
	return
}

func getNodeBandwidth(args framwork.Arguments, topologyNodeList []string) error {
	argv, ok := args[XPUTopologyNodeBandwidth]
	if !ok {
		return errors.New("XPUTopologyNodeBandwidth not exist")
	}
	value, ok := argv.(string)
	if !ok {
		return errors.New("XPUTopologyNodeBandwidth is not string")
	}
	matrix := strings.Split(value, util.Semicolon)
	if len(matrix) != len(topologyNodeList) {
		return errors.New("length of node bandwidth matrix is different from length of node list")
	}
	nodeBandWidth, err := util.ConverMatrix2Map(matrix, topologyNodeList)
	if err != nil {
		util.XPUTopologyNodeBandwidth = nil
		return err
	}
	util.XPUTopologyNodeBandwidth = nodeBandWidth
	klog.V(util.LogInfoLevel).Infof("XPUTopologyNodeBandwidth: +%v", util.XPUTopologyNodeBandwidth)
	return nil
}

func addJobValidFn(ssn *framwork.Session, xp *huaweiXPUPlugin) {
	// check job npu resource, if illegal return failed
	ssn.AddJobValidFn(xp.Name(), func(obj interface{}) *api.ValidateResult {
		return xp.Scheduler.JobValid(obj)
	})
}

func addPredicateFn(ssn *framwork.Session, xp *huaweiXPUPlugin) {
	// if node not meet the task require, the task will be failed. so need to intercept in advance
	ssn.AddPredicateFn(xp.Name(), func(taskInfo *api.TaskInfo, nodeInfo *api.NodeInfo) ([]*api.Status, error) {
		if err != nil {
			xp.Scheduler.Jobs[taskInfo.Job].Lock()
			xp.Scheduler.Jobs[taskInfo.Job].Reason[err.Error()] += nodeInfo.Name + " "
			xp.Scheduler.Jobs[taskInfo.Job].Unlock()
		}
		return predicateStatus, err
	})
}
