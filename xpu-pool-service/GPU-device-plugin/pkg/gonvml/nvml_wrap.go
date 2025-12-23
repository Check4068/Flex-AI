package gonvml

import (
	"errors"
	"fmt"
	"unsafe"
)

// #cgo CFLAGS: -I /usr/local/cuda/targets/x86_64-linux/include -DNVML_NO_UNVERSIONED_FUNC_DEFS=1 -fstack-protector-all
// #cgo LDFLAGS: -L/usr/local/cuda/lib64 -lnvidia-ml -ldl
/*
#include <stddef.h>
#include "dlfcn.h"
#include <stdlib.h>
#include <stdio.h>
#include "nvml.h"

const char* defaultNvmlLibraryName = "libnvidia-ml.so.1";
void *nvmlHandle;

typedef nvmlReturn_t (*NvmlInitFunc)(void);
typedef nvmlReturn_t (*NvmlInitWithFlagsFunc)(unsigned int flags);
typedef nvmlReturn_t (*NvmlShutdownFunc)(void);
typedef const char* (*NvmlErrorStringFunc)(nvmlReturn_t result);
typedef nvmlReturn_t (*NvmlDeviceGetCountFunc)(unsigned int *deviceCount);
typedef nvmlReturn_t (*NvmlDeviceGetHandleByIndexFunc)(unsigned int index, nvmlDevice_t *device);
typedef nvmlReturn_t (*NvmlDeviceGetHandleByUUIDFunc)(const char *uuid, nvmlDevice_t *device);
typedef nvmlReturn_t (*NvmlDeviceGetMemoryInfoV2Func)(nvmlDevice_t device, nvmlMemory_v2_t *memory);
typedef nvmlReturn_t (*NvmlDeviceGetNameFunc)(nvmlDevice_t device, char *name, unsigned int length);
typedef nvmlReturn_t (*NvmlDeviceGetUUIDFunc)(nvmlDevice_t device, char *uuid, unsigned int length);
typedef nvmlReturn_t (*NvmlDeviceGetIndexFunc)(nvmlDevice_t device, unsigned int *index);
typedef nvmlReturn_t (*NvmlDeviceRegisterEventsFunc)(nvmlDevice_t device, unsigned long long eventTypes, nvmlEventSet_t set);
typedef nvmlReturn_t (*NvmlEventSetCreateFunc)(nvmlEventSet_t *set);
typedef nvmlReturn_t (*NvmlEventSetWaitFunc)(nvmlEventSet_t set, nvmlEventData_t * data, unsigned int timeoutms);
typedef nvmlReturn_t (*NvmlEventSetFreeFunc)(nvmlEventSet_t set);
typedef nvmlReturn_t (*NvmlDeviceGetUtilizationRatesFunc)(nvmlDevice_t device, nvmlUtilization_t *utilization);
typedef nvmlReturn_t (*NvmlDeviceGetComputeRunningProcessesFunc)(nvmlDevice_t device, unsigned int* infoSize, nvmlProcessInfo_v1_t* infos);
typedef nvmlReturn_t (*NvmlDeviceGetProcessUtilizationFunc)(nvmlDevice_t device, nvmlProcessUtilizationSample_t* samples, unsigned int* sampleSize, unsigned long long timestamp);
typedef nvmlReturn_t (*NvmlDeviceGetMultiGpuBoardFunc)(nvmlDevice_t device, unsigned int *multiGpuBool);
typedef nvmlReturn_t (*NvmlDeviceGetTopologyCommonAncestorFunc)(nvmlDevice_t device1, nvmlDevice_t device2, nvmlGpuTopologyLevel_t *pathInfo);
typedef nvmlReturn_t (*NvmlDeviceGetTopologyNearestGpusFunc)(nvmlDevice_t device, nvmlGpuTopologyLevel_t level, unsigned int *count, nvmlDevice_t *deviceArray);
typedef nvmlReturn_t (*NvmlSystemGetDriverVersionFunc)(char *version, unsigned int length);
typedef nvmlReturn_t (*NvmlSystemGetCudaDriverVersionFunc)(int *cudaDriverVersion);
typedef nvmlReturn_t (*NvmlDeviceGetTemperatureFunc)(nvmlDevice_t device, nvmlTemperatureSensors_t sensorType, unsigned int *temp);
typedef nvmlReturn_t (*NvmlDeviceGetPowerUsageFunc)(nvmlDevice_t device, unsigned int *power);

NvmlInitFunc nvmlInitFunc = NULL;
NvmlInitWithFlagsFunc nvmlInitWithFlagsFunc = NULL;
NvmlShutdownFunc nvmlShutdownFunc = NULL;
NvmlErrorStringFunc nvmlErrorStringFunc = NULL;
NvmlDeviceGetCountFunc nvmlDeviceGetCountFunc = NULL;
NvmlDeviceGetHandleByIndexFunc nvmlDeviceGetHandleByIndexFunc = NULL;
NvmlDeviceGetHandleByUUIDFunc nvmlDeviceGetHandleByUUIDFunc = NULL;
NvmlDeviceGetMemoryInfoV2Func nvmlDeviceGetMemoryInfoV2Func = NULL;
NvmlDeviceGetNameFunc nvmlDeviceGetNameFunc = NULL;
NvmlDeviceGetUUIDFunc nvmlDeviceGetUUIDFunc = NULL;
NvmlDeviceGetIndexFunc nvmlDeviceGetIndexFunc = NULL;
NvmlDeviceRegisterEventsFunc nvmlDeviceRegisterEventsFunc = NULL;
NvmlEventSetCreateFunc nvmlEventSetCreateFunc = NULL;
NvmlEventSetWaitFunc nvmlEventSetWaitFunc = NULL;
NvmlEventSetFreeFunc nvmlEventSetFreeFunc = NULL;
NvmlDeviceGetUtilizationRatesFunc nvmlDeviceGetUtilizationRatesFunc = NULL;
NvmlDeviceGetComputeRunningProcessesFunc nvmlDeviceGetComputeRunningProcessesFunc = NULL;
NvmlDeviceGetProcessUtilizationFunc nvmlDeviceGetProcessUtilizationFunc = NULL;
NvmlDeviceGetMultiGpuBoardFunc nvmlDeviceGetMultiGpuBoardFunc = NULL;
NvmlDeviceGetTopologyCommonAncestorFunc nvmlDeviceGetTopologyCommonAncestorFunc = NULL;
NvmlDeviceGetTopologyNearestGpusFunc nvmlDeviceGetTopologyNearestGpusFunc = NULL;
NvmlSystemGetDriverVersionFunc nvmlSystemGetDriverVersionFunc = NULL;
NvmlSystemGetCudaDriverVersionFunc nvmlSystemGetCudaDriverVersionFunc = NULL;
NvmlDeviceGetTemperatureFunc nvmlDeviceGetTemperatureFunc = NULL;
NvmlDeviceGetPowerUsageFunc nvmlDeviceGetPowerUsageFunc = NULL;

nvmlReturn_t nvmlInit(void) {
	return (nvmlInitFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlInitFunc();
}

nvmlReturn_t nvmlInitWithFlags(unsigned int flags) {
	return (nvmlInitWithFlagsFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlInitWithFlagsFunc(flags);
}

nvmlReturn_t nvmlShutdown() {
	return (nvmlShutdownFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlShutdownFunc();
}

const char* nvmlErrorString(nvmlReturn_t result) {
	return (nvmlErrorStringFunc == NULL) ? "Function not found" : nvmlErrorStringFunc(result);
}

nvmlReturn_t nvmlSystemGetDriverVersion(char *version, unsigned int length) {
	return (nvmlSystemGetDriverVersionFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlSystemGetDriverVersionFunc(version, length);
}

nvmlReturn_t nvmlSystemGetCudaDriverVersion(int *cudaDriverVersion) {
	return (nvmlSystemGetCudaDriverVersionFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlSystemGetCudaDriverVersionFunc(cudaDriverVersion);
}

nvmlReturn_t nvmlDeviceGetTemperature(nvmlDevice_t device, nvmlTemperatureSensors_t sensorType, unsigned int *temp) {
	return (nvmlDeviceGetTemperatureFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetTemperatureFunc(device, sensorType, temp);
}

nvmlReturn_t nvmlDeviceGetPowerUsage(nvmlDevice_t device, unsigned int *power) {
	return (nvmlDeviceGetPowerUsageFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetPowerUsageFunc(device, power);
}

nvmlReturn_t nvmlDeviceGetCount(unsigned int *deviceCount) {
    return (nvmlDeviceGetCountFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetCountFunc(deviceCount);
}

nvmlReturn_t nvmlDeviceGetHandleByIndex(unsigned int index, nvmlDevice_t *device) {
    return (nvmlDeviceGetHandleByIndexFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetHandleByIndexFunc(index, device);
}

nvmlReturn_t nvmlDeviceGetHandleByUUIDHook(const char *uuid, nvmlDevice_t *device) {
    return (nvmlDeviceGetHandleByUUIDFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetHandleByUUIDFunc(uuid, device);
}

nvmlReturn_t nvmlDeviceGetMemoryInfo_v2Hook(nvmlDevice_t device, nvmlMemory_v2_t *memory) {
    return (nvmlDeviceGetMemoryInfoV2Func == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetMemoryInfoV2Func(device, memory);
}

nvmlReturn_t nvmlDeviceGetName(nvmlDevice_t device, char *name, unsigned int length) {
    return (nvmlDeviceGetNameFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetNameFunc(device, name, length);
}

nvmlReturn_t nvmlDeviceGetUUID(nvmlDevice_t device, char *uuid, unsigned int length) {
    return (nvmlDeviceGetUUIDFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetUUIDFunc(device, uuid, length);
}

nvmlReturn_t nvmlDeviceGetIndex(nvmlDevice_t device, unsigned int *index) {
    return (nvmlDeviceGetIndexFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetIndexFunc(device, index);
}

nvmlReturn_t nvmlDeviceRegisterEvents(nvmlDevice_t device, unsigned long long eventTypes, nvmlEventSet_t set) {
    return (nvmlDeviceRegisterEventsFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceRegisterEventsFunc(device, eventTypes, set);
}

nvmlReturn_t nvmlEventSetCreate(nvmlEventSet_t *set) {
    return (nvmlEventSetCreateFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlEventSetCreateFunc(set);
}

nvmlReturn_t nvmlEventSetWait(nvmlEventSet_t set, nvmlEventData_t *data, unsigned int timeoutms) {
    return (nvmlEventSetWaitFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlEventSetWaitFunc(set, data, timeoutms);
}

nvmlReturn_t nvmlEventSetFree(nvmlEventSet_t set) {
    return (nvmlEventSetFreeFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlEventSetFreeFunc(set);
}

nvmlReturn_t nvmlDeviceGetUtilizationRates(nvmlDevice_t device, nvmlUtilization_t *utilization) {
    return (nvmlDeviceGetUtilizationRatesFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetUtilizationRatesFunc(device, utilization);
}

nvmlReturn_t nvmlDeviceGetComputeRunningProcesses_v1(nvmlDevice_t device, unsigned int *infoSize, nvmlProcessInfo_v1_t *infos) {
    return (nvmlDeviceGetComputeRunningProcessesFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetComputeRunningProcessesFunc(device, infoSize, infos);
}

nvmlReturn_t nvmlDeviceGetProcessUtilization(nvmlDevice_t device, nvmlProcessUtilizationSample_t *samples, unsigned int *sampleSize, unsigned long long timestamp) {
    return (nvmlDeviceGetProcessUtilizationFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetProcessUtilizationFunc(device, samples, sampleSize, timestamp);
}

nvmlReturn_t nvmlDeviceGetMultiGpuBoard(nvmlDevice_t device, unsigned int *multiGpuBool) {
	return (nvmlDeviceGetMultiGpuBoardFunc ==NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetMultiGpuBoardFunc(device, multiGpuBool);
}

nvmlReturn_t nvmlDeviceGetTopologyCommonAncestor(nvmlDevice_t device1, nvmlDevice_t device2, nvmlGpuTopologyLevel_t *pathInfo) {
	return (nvmlDeviceGetTopologyCommonAncestorFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetTopologyCommonAncestorFunc(device1, device2, pathInfo);
}

nvmlReturn_t nvmlDeviceGetTopologyNearestGpus(nvmlDevice_t device, nvmlGpuTopologyLevel_t level, unsigned int *count, nvmlDevice_t *deviceArray) {
	return (nvmlDeviceGetTopologyNearestGpusFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetTopologyNearestGpusFunc(device, level, count, deviceArray);
}

static void loadSymbol(const char *symbolName, void **symbolPtr) {
  *symbolPtr = dlsym(nvmlHandle, symbolName);
  if (!*symbolPtr) {
    fprintf(stderr, "Failed to load symbol %s\n", symbolName);
  }
}

nvmlReturn_t loadDlFunction(void) {
  nvmlHandle = dlopen("libnvidia-ml.so.1", RTLD_LAZY);
  if (nvmlHandle == NULL) {
    return NVML_ERROR_LIBRARY_NOT_FOUND;
  }
  loadSymbol("nvmlInit_v2", (void**)(&nvmlInitFunc));
  loadSymbol("nvmlInitWithFlags", (void**)(&nvmlInitWithFlagsFunc));
  loadSymbol("nvmlShutdown", (void**)(&nvmlShutdownFunc));
  loadSymbol("nvmlErrorString", (void**)(&nvmlErrorStringFunc));
  loadSymbol("nvmlDeviceGetCount_v2", (void**)(&nvmlDeviceGetCountFunc));
  loadSymbol("nvmlDeviceGetHandleByIndex_v2", (void**)(&nvmlDeviceGetHandleByIndexFunc));
  loadSymbol("nvmlDeviceGetHandleByUUID", (void**)(&nvmlDeviceGetHandleByUUIDFunc));
  loadSymbol("nvmlDeviceGetMemoryInfo_v2", (void**)(&nvmlDeviceGetMemoryInfoV2Func));
  loadSymbol("nvmlDeviceGetName", (void**)(&nvmlDeviceGetNameFunc));
  loadSymbol("nvmlDeviceGetUUID", (void**)(&nvmlDeviceGetUUIDFunc));
  loadSymbol("nvmlDeviceGetIndex", (void**)(&nvmlDeviceGetIndexFunc));
  loadSymbol("nvmlDeviceRegisterEvents", (void**)(&nvmlDeviceRegisterEventsFunc));
  loadSymbol("nvmlEventSetCreate", (void**)(&nvmlEventSetCreateFunc));
  loadSymbol("nvmlEventSetWait_v2", (void**)(&nvmlEventSetWaitFunc));
  loadSymbol("nvmlEventSetFree", (void**)(&nvmlEventSetFreeFunc));
  loadSymbol("nvmlDeviceGetUtilizationRates", (void**)(&nvmlDeviceGetUtilizationRatesFunc));
  loadSymbol("nvmlDeviceGetComputeRunningProcesses", (void**)(&nvmlDeviceGetComputeRunningProcessesFunc));
  loadSymbol("nvmlDeviceGetProcessUtilization", (void**)(&nvmlDeviceGetProcessUtilizationFunc));
  loadSymbol("nvmlDeviceGetMultiGpuBoard", (void**)(&nvmlDeviceGetMultiGpuBoardFunc));
  loadSymbol("nvmlDeviceGetTopologyCommonAncestor", (void**)(&nvmlDeviceGetTopologyCommonAncestorFunc));
  loadSymbol("nvmlDeviceGetTopologyNearestGpus", (void**)(&nvmlDeviceGetTopologyNearestGpusFunc));
  loadSymbol("nvmlSystemGetDriverVersion", (void**)(&nvmlSystemGetDriverVersionFunc));
  loadSymbol("nvmlSystemGetCudaDriverVersion", (void**)(&nvmlSystemGetCudaDriverVersionFunc));
  loadSymbol("nvmlDeviceGetTemperature", (void**)(&nvmlDeviceGetTemperatureFunc));
  loadSymbol("nvmlDeviceGetPowerUsage", (void**)(&nvmlDeviceGetPowerUsageFunc));

  fprintf(stdout, "Loaded nvml library\n");
  return NVML_SUCCESS;
}

nvmlReturn_t unloadDlFunction(void) {
  if (nvmlHandle == NULL) {
    return NVML_SUCCESS;
  }
  if (dlclose(nvmlHandle) == 0) {
    nvmlHandle = NULL;
    return NVML_SUCCESS;
  }
  return NVML_ERROR_UNKNOWN;
}
*/

import "C"

var errNvmlDlLoaded = errors.New("could not load nvml library")

type nvmlDevice struct {
	Handle C.nvmlDevice_t
}

type nvmlEventSet struct {
	Handle C.nvmlEventSet_t
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
	cFlags, _ := (C.uint)(Flags), cgoAllocsUnknown
	return NvmlRetType(C.nvmlInitWithFlags(cFlags))
}

func nvmlDeviceGetCountWrapper(DeviceCount *uint32) NvmlRetType {
	cDeviceCount, _ := (*C.uint)(unsafe.Pointer(DeviceCount)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetCount(cDeviceCount))
}

func nvmlSystemGetDriverVersionWrapper(Version *byte, Length uint32) NvmlRetType {
	cVersion, _ := (*C.char)(unsafe.Pointer(Version)), cgoAllocsUnknown
	cLength, _ := (C.uint)(Length), cgoAllocsUnknown
	return (NvmlRetType)(C.nvmlSystemGetDriverVersion(cVersion, cLength))
}

func nvmlSystemGetCudaDriverVersionWrapper(CudaDriverVersion *int32) NvmlRetType {
	cCudaDriverVersion, _ := (*C.int)(unsafe.Pointer(CudaDriverVersion)), cgoAllocsUnknown
	return (NvmlRetType)(C.nvmlSystemGetCudaDriverVersion(cCudaDriverVersion))
}

func nvmlDeviceGetHandleByIndexWrapper(Index uint32, Device *nvmlDevice) NvmlRetType {
	cIndex, _ := (C.uint)(Index), cgoAllocsUnknown
	cDevice, _ := (*C.nvmlDevice_t)(unsafe.Pointer(Device)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetHandleByIndex(cIndex, cDevice))
}

func nvmlDeviceGetHandleByUUIDWrapper(Uuid string, Device *nvmlDevice) NvmlRetType {
	cUuid, _ := unPackPCharString(Uuid)
	cDevice, _ := (*C.nvmlDevice_t)(unsafe.Pointer(Device)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetHandleByUUIDHook(cUuid, cDevice))
}

func nvmlDeviceGetMemoryInfoWrapper(nvmlDevice nvmlDevice, memory *MemoryV2) NvmlRetType {
	cnvmlDevice, _ := *(*C.nvmlDevice_t)(unsafe.Pointer(&nvmlDevice)), cgoAllocsUnknown
	cmemory, _ := (*C.nvmlMemory_v2_t)(unsafe.Pointer(memory)), cgoAllocsUnknown
	return NvmlRetType(C.nvmlDeviceGetMemoryInfo_v2Hook(cnvmlDevice, cmemory))
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
	cData, _ := (*C.nvmlEventData_t)(unsafe.Pointer(Data)), cgoAllocsUnknown
	cTimeouts, _ := (C.uint)(Timeouts), cgoAllocsUnknown
	return NvmlRetType(C.nvmlEventSetWait(cSet, cData, cTimeouts))
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
