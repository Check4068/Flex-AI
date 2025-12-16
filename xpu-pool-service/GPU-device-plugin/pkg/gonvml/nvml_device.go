/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2025. All rights reserved.
 */

// Package gonvml implements accessing the NVML library using the go
package gonvml

import (
	"C"
	"fmt"
	"reflect"
)

const deviceGetMemInfoVersion = 2

var pidMaxSize uint32 = 1024

func (device NvmlDevice) GetMemoryInfoV2() (MemoryV2, NvmlRetType) {
	var memory MemoryV2
	memory.Version = structVersion(memory, deviceGetMemInfoVersion)
	ret := nvmlDeviceGetMemoryInfoWrapper(device, &memory)
	return memory, ret
}

func (device NvmlDevice) GetName() (string, NvmlRetType) {
	name := make([]byte, DeviceNameV2BufferSize)
	ret := libnvml.DeviceGetNameV2(device, name[:], DeviceNameV2BufferSize)
	return string(name[:clen(name)]), ret
}

func (device NvmlDevice) RegisterEvent(eventTypes uint64, set EventSet) NvmlRetType {
	return libnvml.DeviceRegisterEventWrapper(device, eventTypes, set.(*nvmlEventSet))
}

func (device NvmlDevice) GetUUID() (string, NvmlRetType) {
	uuid := make([]byte, DeviceUUIDBufferSize)
	ret := libnvml.DeviceGetUUIDWrapper(device, uuid[:], DeviceUUIDBufferSize)
	return string(uuid[:clen(uuid)]), ret
}

func (device NvmlDevice) GetIndex() (int, NvmlRetType) {
	var index uint32
	ret := libnvml.DeviceGetIndexWrapper(device, &index)
	return int(index), ret
}

func (device NvmlDevice) GetUtilizationRates() (Utilization, NvmlRetType) {
	var utilization Utilization
	ret := libnvml.DeviceGetUtilizationRatesWrapper(device, utilization)
	return utilization, ret
}

func (device NvmlDevice) GetComputeRunningProcesses() ([]ProcessInfoV2, NvmlRetType) {
	var infoSize uint32 = pidMaxSize
	var infos = make([]ProcessInfoV2, infoSize)
	ret := libnvml.DeviceGetComputeRunningProcessesWrapper(device, &infoSize, infos[:])
	return infos, ret
}

func (device NvmlDevice) GetProcessUtilization(timestamp uint64) ([]ProcessUtilizationSample, NvmlRetType) {
	var sampleSize uint32 = pidMaxSize
	var samples = make([]ProcessUtilizationSample, sampleSize)
	ret := libnvml.DeviceGetProcessUtilizationWrapper(device, samples[:], sampleSize, timestamp)
	return samples, ret
}

func (device NvmlDevice) GetMultiGpuBoard() (int, NvmlRetType) {
	var multiGpuBoard uint32
	ret := libnvml.DeviceGetMultiGpuBoardWrapper(device, &multiGpuBoard)
	return int(multiGpuBoard), ret
}

func (device1 NvmlDevice) GetTopologyCommonAncestor(device2 Device, level GpuTopologyLevel) (Device, NvmlRetType) {
	var pathInfo Device
	ret := libnvml.DeviceGetTopologyCommonAncestorWrapper(device1, hwDeviceHandle(device2), &pathInfo)
	return pathInfo, ret
}

func nvmlDeviceHandle(device Device) NvmlDevice {
	var helper func(val reflect.Value) NvmlDevice
	helper = func(val reflect.Value) NvmlDevice {
		if val.Kind() == reflect.Interface {
			val = val.Elem()
		}
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Type() == reflect.TypeOf(NvmlDevice{}) {
			return val.Interface().(NvmlDevice)
		}
		if val.Kind() != reflect.Struct {
			panic(fmt.Errorf("unable to convert non-struct type %v to NvmlDevice", val.Kind()))
		}
		for i := 0; i < val.Type().NumField(); i++ {
			if val.Field(i).Anonymous {
				continue
			}
			if !val.Field(i).Type().Implements(reflect.TypeOf((*Device)(nil)).Elem()) {
				continue
			}
			return helper(val.Field(i))
		}
		panic(fmt.Errorf("unable to convert %T to NvmlDevice", d))
	}
	return helper(reflect.ValueOf(d))
}

func (device NvmlDevice) GetTopologyNearestGpus(level GpuTopologyLevel) ([]Device, NvmlRetType) {
	var count uint32
	ret := libnvml.DeviceGetTopologyNearestGpusWrapper(device, level, &count, nil)
	if ret != Success {
		return []Device{}, ret
	}
	if count == 0 {
		return []Device{}, ret
	}
	deviceArray := make([]NvmlDevice, count)
	ret = libnvml.DeviceGetTopologyNearestGpusWrapper(device, level, &count, &deviceArray[0])
	return convertDevices(deviceArray), ret
}

func (device NvmlDevice) GetTemperature(temperatureGpu NvmlTemperatureSensors) (uint32, NvmlRetType) {
	var temp uint32
	ret := libnvml.DeviceGetTemperatureWrapper(device, temperatureGpu, &temp)
	return temp, ret
}

func (device NvmlDevice) GetPowerUsage() (uint32, NvmlRetType) {
	var power uint32
	ret := libnvml.DeviceGetPowerUsageWrapper(device, &power)
	return power, ret
}