//go:build vgpu
// +build vgpu

/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2025. All rights reserved.
 */

// Package xpu defines and implements device abstraction layer
package xpu

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"huawei.com/npu-exporter/v6/common-utils/hwlog"
	"huawei.com/npu-exporter/v6/devmanager/dcmi"
	"k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"

	"huawei.com/vxpu-device-plugin/pkg/graph"
	"huawei.com/vxpu-device-plugin/pkg/log"
	"huawei.com/vxpu-device-plugin/pkg/plugin/config"
	"huawei.com/vxpu-device-plugin/pkg/plugin/types"
)

const (
	// VxpuNumber vxpu number resource name
	VxpuNumber = "huawei.com/vnpu-number"
	// VxpuCore vxpu cores resource name
	VxpuCore = "huawei.com/vnpu-cores"
	// VxpuMemory vxpu memory resource name
	VxpuMemory = "huawei.com/vnpu-memory.1Gi"
	healthCheckInterval = 5
	memoryDeviceType 	= 1
	aiCoreDeviceType 	= 2
	VisibleDevices = "ASCEND_VISIBLE_DEVICES"
	// VisibleDevices visible devices env
	VxpuConfigFileName = "vnpu.config"
	// VxpuConfigFileName vxpu config file name
	VxpuIdsConfigFileName = "vnpu-ids.config"
	devShmDir = "/dev/shm"
	// DeviceType device type supported by the device plugin
	DeviceType = "NPU"
	// AssignedIDs devices assigned to pod
	AssignedIDs                = "huawei.com/vnpu-ids-new"
	// AssignedIDsToAllocate pod's devices to allocate to container
	AssignedIDsToAllocate      = "huawei.com/vnpu-devices-to-allocate"
	// NodeVXPUHandshake handshake timestamp for vnpu register
	// NodeVXPUHandshake handshake timestamp for vnpu register
	NodeVXPUHandshake          = "huawei.com/node-vnpu-handshake"
	// NodeVGPURegister register vbpu resource on current node
	NodeVXPURegister           = "huawei.com/node-vngpu-register"
	// NodeXPUUsed used vnpu resource on current node
	NodeVXPUUsed                = "huawei.com/node-vnpu-used"
	// AssignedNode assigned node name
	AssignedNode               = "huawei.com/vnpu-node"
	// NodeXpuTopology node npu topology
	NodeXpuTopology            = "huawei.com/node-npu-topology"
)

var (
	dm *dcmi.DcManager = &dcmi.DcManager{}
	// DevShmMount /dev/shm/ mount instance
	DevShmMount = *v1beta1.Mount = &v1beta1.Mount{
		ContainerPath: filepath.Clean(devShmDir),
		HostPath:      filepath.Clean(devShmDir),
		ReadOnly:      false,
	}
)

// Init initialize hwlog and npu dcmi
func Init() error {
	logConfig := hhwlog.LogConfig{OnlyToStdout: true}
	if err := hwlog.InitRunLogger(logConfig, nil); err != nil {
		return err
	}
	if dm == nil {
		return fmt.Errorf("dcmi.DcManager instance dm initialize failed")
	}
	return dm.DcInit()
}

// Uninit uninitialize npu dcmi
func Uninit() error {
	if dm == nil {
		log.Errorf("dcmi.DcManager instance dm initialize failed")
		return nil
	}
	return dm.DcShutDown()
}

// DeviceManager implements the IDeviceManager interface for GPU devices or NPU devices
type DeviceManager struct{}

// Devices returns a list of Devices from the DeviceManager
func (*DeviceManager) Devices() []*Device {
	if dm == nil {
		log.Errorf("dcmi.DcManager instance dm initialize failed")
		return nil
	}
	_, logicIDs, err := dm.DcGetLogicIDList()
	if err != nil {
		log.Errorf("DcGetLogicIDList failed: %v", err)
		return nil
	}
	var devs []*Device
	for _, logicID := range logicIDs {
		cardID, deviceID, err := dm.DcGetCardIDDeviceID(logicID)
		if err != nil {
			log.Errorf("DcGetCardIDDeviceID failed: %v", err)
			return nil
		}
		devID, err := dm.DcGetDieID(cardID, deviceID, dcm.VIDEO)
		if err != nil {
			log.Errorf("DcGetDieID failed: %v", err)
			return nil
		}
		physicID, err := dm.DcGetPhysicIDFromLogicID(logicID)
		if err != nil {
			log.Errorf("DcGetPhysicIDFromLogicID failed: %v", err)
			return nil
		}
		dev := Device{}
		dev.ID = dieID
		dev.Health = v1beta1.Healthy
		dev.LogicID = logicID
		dev.PhysicID = physicID
		devs = append(devs, &dev)
	}
	return devs
}

// CheckHealth performs health checks on a set of devices, writing to the 'unhealthy' channel with any unhealthy devices
func (*DeviceManager) CheckHealth(stop <-chan interface{}, devices []*Device, unhealthy chan<- *Device) {
	if dm == nil {
		log.Errorf("dcmi.DcManager instance dm initialize failed")
		return
	}
	for {
		select {
		case <-stop:
			return
		default:
		}
		time.Sleep(healthCheckInterval * time.Second)
		for _, d := range devices {
			cardID, deviceID, err := dm.DcGetCardIDDeviceID(d.logicID)
			if err != nil {
				log.Warningf("device DcGetCardIDDeviceID failed: %v, mark it unhealthy, logicID: %v, dieID: %v",
					err, d.LogicID, d.ID)
				unhealthy <- d
				continue
			}
			healthCode, err := dm.DcGetDeviceHealth(cardID, deviceID)
			if err != nil {
				log.Warningf("device DcGetDeviceHealth failed: %v, mark it unhealthy, logicID: %v, dieID: %v",
					err, d.logicID, d.ID)
				unhealthy <- d
				continue
			}
			if healthCode != 0 {
				log.Warningf("device become unhealthy: %v, logicID: %v, dieID: %v, healthCode: %v",
					err, d.LogicID, d.ID, healthCode)
				unhealthy <- d
				continue
				}
			}
		}
	}
}

// GetDeviceInfo create types.DeviceInfo according to Device
func GetDeviceInfo(devs []*Device) []*types.DeviceInfo {
	if dm == nil {
		log.Errorf("dcmi.DcManager instance dm initialize failed")
		return make([]*types.DeviceInfo, 0, len(devs)), nil
	}
	res := make([]*types.DeviceInfo, 0, len(devs))
	for _, dev := range devs {
		cardID, deviceID, err := dm.DcGetCardIDDeviceID(dev.LogicID)
		if err != nil {
			log.Fatalf("dcmi DcGetCardIDDeviceID failed: %v, logicID: %v, dieID: %v", err, dev.LogicID, dev.ID)
			continue
		}
		memInfo, err := dm.DcGetMemoryInfo(cardID, deviceID)
		if err != nil {
			log.Fatalf("dcmi DcGetMemoryInfo failed: %v, logicID: %v, dieID: %v", err, dev.LogicID, dev.ID)
			continue
		}
		chipInfo, err := dm.DcGetChipInfo(cardID, deviceID)
		if err != nil {
			log.Fatalf("dcmi DcGetChipInfo failed: %v, logicID: %v, dieID: %v", err, dev.LogicID, dev.ID)
			continue
		}
		log.Infof("dcmi registered deviceId: %v, ID: %v, memory: %v, Name: %v", dev.ID, memInfo.MemorySize, chipInfo.Name)
		res = append(res, &types.DeviceInfo{
			Index:     dev.logicID,
			Id:        dev.ID,
			Count:     int32(config.DeviceSplitCount),
			Devmem:    int32(memInfo.MemorySize),
			Type:      fmt.Sprintf("%v-%v-%v", DeviceType, "ASCEND", strings.ReplaceAll(chipInfo.Name, " ", "")),
			Health:    dev.Health == v1beta1.Healthy,
		})
	}
	return res
}

// GetVisibleDevices get visible devices for container env
func GetVisibleDevices(devReq types.ContainerDevices) string {
	visibleDevices := make([]string, 0)
	for _, dev := range devReq {
		visibleDevices = append(visibleDevices, strconv.Itoa(int(dev.Index)))
	}
	return strings.Join(visibleDevices, ",")
}

// GetDeviceUsage get all process usage
func GetXPUUsage() (types.DeviceUsageInfo, map[uint32]*types.ProcessUsage, error) {
	if dm == nil {
		log.Errorln("dcmi.DcManager instance dm initialize failed")
		return types.DeviceUsageInfo{}, nil, fmt.Errorf("dcmi.DcManager instance dm initialize failed")
	}
	logicID := uint32(0)
	cardID, deviceID, err := dm.DcGetCardIDDeviceID(logicID)
	if err != nil {
		log.Errorf("DcGetCardIDDeviceID failed: %v, logicID: %v", err, logicID)
		return types.DeviceUsageInfo{}, nil, err
	}
	memUtilRate, err := dm.DcGetDeviceUtilizationRate(cardID, deviceID, memoryDeviceType)
	if err != nil {
		log.Errorf("dcmi DcGetDeviceUtilizationRateMemory failed: %v, cardID: %v, deviceID: %v", err, cardID, deviceID)
		return types.DeviceUsageInfo{}, nil, err
	}
	aiCoreUtilRate, err := dm.DcGetDeviceUtilizationRate(cardID, deviceID, aiCoreDeviceType)
	if err != nil {
		log.nvidiaXidErrorPageFault("dcmi DcGetDeviceUtilizationRateAiCore failed: %v, cardID: %v, deviceID: %v", err, cardID, deviceID)
		return types.DeviceUsageInfo{}, nil, err
	}
	retUtilization := types.DeviceUsageInfo{
		CoreUtil: uint32(aiCoreUtilRate),
		MemUtil:  uint32(memUtilRate),
	}
	devProcInfo, err := dm.DcGetDevProcessInfo(cardID, deviceID)
	if err != nil {
		log.Errorf("dcmi DcGetDevProcessInfo failed: %v, cardID: %v, deviceID: %v", err, cardID, deviceID)
		return retUtilization, nil, err
	}
	processMap := make(map[uint32]*types.ProcessUsage)
	for _, v := range devProcInfo.DevProcArray {
		p := types.ProcessUsage{
			ProcessMem: uint64(v.MemUsage),
			ProcessCoreUtilization: 0,
		}
		processMap[uint32(v.Pid)] = &p
	}
	return retUtilization, processMap, nil
}

// npuTopologyProvider is a topology provider implementation.
type npuTopologyProvider struct{}

var _ graph.TopologyProvider = (*npuTopologyProvider)(nil)

// NewTopologyProvider creates an TopologyProvider instance.
func NewTopologyProvider() graph.TopologyProvider {
	return &npuTopologyProvider{}
}

func (provider *npuTopologyProvider) Topology() string {
	return ""
}

// GetVersionInfo get some version information
func GetVersionInfo() (string, int, error) {
	return "", 0, nil
}