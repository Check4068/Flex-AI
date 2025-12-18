package gonvml

import (
	"errors"
	"fmt"
	"unsafe"
)

import "C"

var errNvmlDlLoaded = errors.New("could not load nvml library")

type nvmlDevice struct {
	handle C.nvmlDevice_t
}

type nvmlEventSet struct {
	handle C.nvmlEventSet_t
}

func errorString(ret C.nvmlReturn_t) error {
	if ret == C.NVML_SUCCESS {
		return nil
	}
	if ret == C.NVML_ERROR_LIBRARY_NOT_FOUND || C.nvmlHandle == nil {
		return errNvmlDlLoaded
	}
	err := fmt.Errorf("nvml error %d: %s", ret, C.GoString(C.nvmlErrorString(ret)))
	return err
}

var cgoAllocsUnknown = new(struct{})

type stringHeader struct {
	Data unsafe.Pointer
	Len  int
}

func unPackPCharString(str string) (*C.char, *struct{}) {
	h := (*stringHeader)(unsafe.Pointer(&str))
	return (*C.char)(h.Data), cgoAllocsUnknown
}

func loadNvmlSo() error {
	err := errorString(C.loadDlFunction())
	return err
}

func unloadNvmlSo() error {
	return errorString(C.unloadDlFunction())
}

func nvmlInitWrapper() NvmlRetType {
	return NvmlRetType(C.nvmlInit())
}

func nvmlShutdownWrapper() NvmlRetType {
	return NvmlRetType(C.nvmlShutdown())
}

func nvmlInitWithFlagsWrapper(Flags uint32) NvmlRetType {
	cFlags, _ := (C.unit)(Flags), cgoAllocsUnknown
	return NvmlRetType(C.nvmlInitWithFlags(cFlags))
}

func nvmlDeviceGetCountWrapper(DeviceCount *uint32) NvmlRetType {
	cDeviceCount, _ := (*C.uint)(unsafe.Pointer(DeviceCount)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetCount(cDeviceCount))
}

func nvmlSystemGetDriverVersionWrapper(Version *byte, Length uint32) NvmlRetType {
	cVersion, _ := (*C.char)(unsafe.Pointer(Version)), cgoAllocsUnknown
	cLength, _ := (C.uint)(Length), cgoAllocsUnknown
	return NvmlRetType(C.nvmlSystemGetDriverVersion(cVersion, cLength))
}

func nvmlSystemGetCudaDriverVersionWrapper(CudaDriverVersion *int32) NvmlRetType {
	cCudaDriverVersion, _ := (*C.int)(unsafe.Pointer(CudaDriverVersion)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlSystemGetCudaDriverVersion(cCudaDriverVersion))
}

func nvmlDeviceGetHandleByIndexWrapper(Index uint32, Device *nvmlDevice) NvmlRetType {
	cIndex, _ := (C.uint)(Index), cgoAllocsUnknown
	cDevice, _ := (*C.nvmlDevice_t)(unsafe.Pointer(Device)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetHandleByIndex(cIndex, cDevice))
}

func nvmlDeviceGetHandleByUUIDWrapper(Uuid string, Device *nvmlDevice) NvmlRetType {
	cUuid, _ := unPackPCharString(Uuid)
	cDevice, _ := (*C.nvmlDevice_t)(unsafe.Pointer(Device)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetHandleByUUID(cUuid, cDevice))
}

func nvmlDeviceGetMemoryInfoWrapper(nvmlDevice nvmlDevice, memory *MemoryV2) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&nvmlDevice)), cgoAllocsUnknown
	cmemory, _ := (*C.nvmlMemory_v2_t)(unsafe.Pointer(memory)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetMemoryInfo_v2hook(cnvmlDevice, cmemory))
}

func nvmlDeviceGetNameWrapper(nvmlDevice nvmlDevice, name *byte, length uint32) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&nvmlDevice)), cgoAllocsUnknown
	cname, _ := (*C.char)(unsafe.Pointer(name)), cgoAllocsUnknown
	clength, _ := (C.uint)(length), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetName(cnvmlDevice, cname, clength))
}

func nvmlDeviceGetUUIDWrapper(device nvmlDevice, uuid *byte, length uint32) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cuuid, _ := (*C.char)(unsafe.Pointer(uuid)), cgoAllocsUnknown
	clength, _ := (C.uint)(length), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetUUID(cnvmlDevice, cuuid, clength))
}

func nvmlDeviceGetIndexWrapper(device nvmlDevice, index *uint32) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cindex, _ := (*C.uint)(unsafe.Pointer(index)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetIndex(cnvmlDevice, cindex))
}

func nvmlDeviceRegisterEventsWrapper(device nvmlDevice, EventsTypes uint64, Set nvmlEventSet) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cEventsTypes, _ := (C.ulonglong)(EventsTypes), cgoAllocsUnknown
	cSet, _ := *(*C.nvmlEventSet_t)(unsafe.Pointer(&Set)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceRegisterEvents(cnvmlDevice, cEventsTypes, cSet))
}

func nvmlEventSetCreateWrapper(Set *nvmlEventSet) NvmlRetType {
	cSet, _ := (*C.nvmlEventSet_t)(unsafe.Pointer(Set)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlEventSetCreate(cSet))
}

func nvmlEventSetWaitWrapper(Set nvmlEventSet, Data *nvmlEventData, Timeouts uint32) NvmlRetType {
	cSet, _ := *(*C.nvmlEventSet_t)(unsafe.Pointer(&Set)), cgoAllocsUnknown
	cTimeouts, _ := (C.uint)(Timeouts), cgoAllocsUnknown
	cData, _ := (*C.nvmlEventData_t)(unsafe.Pointer(Data)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlEventSetWait(cSet, cTimeouts, cData))
}

func nvmlEventSetFreeWrapper(Set nvmlEventSet) NvmlRetType {
	cSet, _ := *(*C.nvmlEventSet_t)(unsafe.Pointer(&Set)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlEventSetFree(cSet))
}

func nvmlDeviceGetUtilizationRatesWrapper(device nvmlDevice, utilization *Utilization) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cutilization, _ := (*C.nvmlUtilization_t)(unsafe.Pointer(utilization)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetUtilizationRates(cnvmlDevice, cutilization))
}

func nvmlDeviceGetComputeRunningProcessesWrapper(device nvmlDevice, infoCount *uint32, infos *ProcessInfoV1) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cinfoCount, _ := (*C.uint)(unsafe.Pointer(infoCount)), cgoAllocsUnknown
	cinfos, _ := (*C.nvmlProcessInfo_v1_t)(unsafe.Pointer(infos)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetComputeRunningProcesses_v1(cnvmlDevice, cinfoCount, cinfos))
}

func nvmlDeviceGetProcessUtilizationWrapper(device nvmlDevice, utilization *ProcessUtilizationSample,
	ProcessCount *uint32, LastSeenTimestamp uint64) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cutilization, _ := (*C.nvmlProcessUtilizationSample_t)(unsafe.Pointer(utilization)), cgoAllocsUnknown
	cProcessCount, _ := (*C.uint)(unsafe.Pointer(ProcessCount)), cgoAllocsUnknown
	clastSeenTimestamp, _ := (C.ulonglong)(LastSeenTimestamp), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetProcessUtilization(cnvmlDevice, cutilization, cProcessCount, clastSeenTimestamp))
}

func nvmlDeviceGetMultiGpuBoardWrapper(device nvmlDevice, MultiGpuBool *uint32) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cisMultiGpuBoard, _ := (*C.uint)(unsafe.Pointer(MultiGpuBool)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetMultiGpuBoard(cnvmlDevice, cisMultiGpuBoard))
}

func nvmlDeviceGetTopologyCommonAncestorWrapper(device1 nvmlDevice, device2 nvmlDevice, PathInfo *GpuTopologyLevel) NvmlRetType {
	cnvmlDevice1, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&device1)), cgoAllocsUnknown
	cnvmlDevice2, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&device2)), cgoAllocsUnknown
	cPathInfo, _ := (*C.nvmlGpuTopologyLevel_t)(unsafe.Pointer(PathInfo)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetTopologyCommonAncestor(cnvmlDevice1, cnvmlDevice2, cPathInfo))
}

func nvmlDeviceGetTopologyNearestGpusWrapper(device nvmlDevice, level GpuTopologyLevel, devicesCount *uint32, devices *nvmlDevice) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&device)), cgoAllocsUnknown
	clevel, _ := (C.nvmlGpuTopologyLevel_t)(level), cgoAllocsUnknown
	cdevicesCount, _ := (*C.uint)(unsafe.Pointer(devicesCount)), cgoAllocsUnknown
	cdevices, _ := (*C.nvmlDevice_t)(unsafe.Pointer(devices)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetTopologyNearestGpus(cnvmlDevice, clevel, cdevicesCount, cdevices))
}

func nvmlDeviceGetTemperatureWrapper(device nvmlDevice, sensorType NvmlTemperatureSensors, temp *uint32) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&device)), cgoAllocsUnknown
	csensorType, _ := *(C.nvmlTemperatureSensors_t)(unsafe.Pointer(&sensorType)), cgoAllocsUnknown
	ctemp, _ := (*C.uint)(unsafe.Pointer(temp)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetTemperature(cnvmlDevice, csensorType, ctemp))
}

func nvmlDeviceGetPowerUsageWrapper(device nvmlDevice, power *uint32) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpower, _ := (*C.uint)(unsafe.Pointer(power)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetPowerUsage(cnvmlDevice, cpower))
}
