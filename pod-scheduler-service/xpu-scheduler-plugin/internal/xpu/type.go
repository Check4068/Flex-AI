package xpu

import (
	"volcano.sh/volcano/pkg/scheduler/plugins/xpu-scheduler-plugin/plugin"
	"volcano.sh/volcano/pkg/scheduler/plugins/xpu-scheduler-plugin/util"
)

var (
	gpuPlugin = &plugin.SchedulerPlugin{
		PluginName:                 util.GpuPluginName,
		VxpuName:                   util.VGpuName,
		VxpuType:                   util.VGpuType,
		VxpuCore:                   util.VGpuCore,
		VxpuMemory:                 util.VGpuMemory,
		Config:                     Config,
		NodeXPURegisterAnno:        util.NodeGPURegisterAnno,
		AssignedXPUsToAllocateAnno: util.AssignedGPUsToAllocateAnnotations,
		AssignedXPUsToNodeAnno:     util.AssignedGPUsToNodeAnnotations,
		AssignedXPUsPodAnno:        util.AssignedGPUsPodAnnotations,
		NodeXPUTopologyAnno:        util.NodeGPUTopologyAnnotation,
		NodeXPUHandshakeAnno:       util.NodeGPUHandshakeAnnotation,
	}

	npuPlugin = &plugin.SchedulerPlugin{
		PluginName:                 util.NPUPluginName,
		VxpuName:                   util.VNPUName,
		VxpuType:                   util.VNPUType,
		VxpuCore:                   util.VNPUCore,
		VxpuMemory:                 util.VNPUMemory,
		Config:                     Config,
		NodeXPURegisterAnno:        util.NodeNPURegisterAnnotation,
		AssignedXPUsToAllocateAnno: util.AssignedNPUsToAllocateAnnotations,
		AssignedXPUsToNodeAnno:     util.AssignedNPUsToNodeAnnotations,
		AssignedXPUsToPodAnno:      util.AssignedNPUsToPodAnnotations,
		NodeXPUTopologyAnno:        util.NodeNPUTopologyAnnotation,
		NodeXPUHandshakeAnno:       util.NodeNPUHandshakeAnnotation,
	}

	Config = &plugin.CommonConfig{}
)
