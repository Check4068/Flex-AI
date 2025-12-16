/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2025. All rights reserved.
 */

// In this file, the cgo feature is used to invoke the NVML library.
// New APIs are supported based on service requirements.

// Package gonvml implements accessing the NVML library using the go

package gonvml

import (
	"errors"
	"fmt"
	"unsafe"
)

/*
#cgo CFLAGS: -I /usr/local/cuda-11.8/targets/x86_64-linux/include -DNO_NVML_UNVERSIONED_FUNC_DEFS=1 -fstack-protector-all
#cgo LDFLAGS: -ldl
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>

static void *nvmlHandle = NULL; // Handle for dynamically loaded libnvidia-ml.so
static const char *nvmlLibraryName = "libnvidia-ml.so.1";

// Define the function prototypes.
typedef nvmlReturn_t (*nvmlInitFunc)(void);
typedef nvmlReturn_t (*nvmlInitWithFlagsFunc)(unsigned int flags);
typedef nvmlReturn_t (*nvmlShutdownFunc)(void);
typedef const char *(*nvmlErrorStringFunc)(nvmlReturn_t result);
typedef nvmlReturn_t (*nvmlDeviceGetCountFunc)(unsigned int *deviceCount);
typedef nvmlReturn_t (*nvmlDeviceGetHandleByIndexFunc)(unsigned int index, nvmlDevice_t *device);
typedef nvmlReturn_t (*nvmlDeviceGetHandleByUUIDFunc)(const char *uuid, nvmlDevice_t *device);
typedef nvmlReturn_t (*nvmlDeviceGetMemoryInfoV2Func)(nvmlDevice_t device, nvmlMemory_v2_t *memory);
typedef nvmlReturn_t (*nvmlDeviceGetNameFunc)(nvmlDevice_t device, char *name, unsigned int length);
typedef nvmlReturn_t (*nvmlDeviceGetUUIDFunc)(nvmlDevice_t device, char *uuid, unsigned int length);
typedef nvmlReturn_t (*nvmlDeviceGetIndexFunc)(nvmlDevice_t device, unsigned int *index);
typedef nvmlReturn_t (*nvmlDeviceRegisterEventsFunc)(nvmlDevice_t device, unsigned long long eventTypes, nvmlEventSet_t set);
typedef nvmlReturn_t (*nvmlEventSetCreateFunc)(nvmlEventSet_t *set);
typedef nvmlReturn_t (*nvmlEventSetWaitFunc)(nvmlEventSet_t set, nvmlEventData_t *data, unsigned int timeouts);
typedef nvmlReturn_t (*nvmlEventSetFreeFunc)(nvmlEventSet_t set);
typedef nvmlReturn_t (*nvmlDeviceGetUtilizationRatesFunc)(nvmlDevice_t device, nvmlUtilization_t *utilization);
typedef nvmlReturn_t (*nvmlDeviceGetComputeRunningProcessesFunc)(nvmlDevice_t device, unsigned int *infoSize, nvmlProcessInfo_v2_t *infos);
typedef nvmlReturn_t (*nvmlDeviceGetProcessUtilizationFunc)(nvmlDevice_t device, nvmlProcessUtilizationSample_t *samples, unsigned int sampleSize, unsigned long long timestamp);
typedef nvmlReturn_t (*nvmlDeviceGetMultiGpuBoardFunc)(nvmlDevice_t device, unsigned int *multiGpuBoard);
typedef nvmlReturn_t (*nvmlDeviceGetTopologyCommonAncestorFunc)(nvmlDevice_t device1, nvmlDevice_t device2, nvmlGpuTopologyLevel_t *pathInfo);
typedef nvmlReturn_t (*nvmlDeviceGetTopologyNearestGpusFunc)(nvmlDevice_t device, nvmlGpuTopologyLevel_t level, unsigned int *count, nvmlDevice_t *deviceArray);
typedef nvmlReturn_t (*nvmlSystemGetDriverVersionFunc)(char *version, unsigned int length);
typedef nvmlReturn_t (*nvmlSystemGetCudaDriverVersionFunc)(int *cudaDriverVersion);
typedef nvmlReturn_t (*nvmlDeviceGetTemperatureFunc)(nvmlDevice_t device, nvmlTemperatureSensors_t sensorType, unsigned int *temp);
typedef nvmlReturn_t (*nvmlDeviceGetPowerUsageFunc)(nvmlDevice_t device, unsigned int *power);

// Function pointers.
static nvmlInitFunc nvmlInitFunc = NULL;
static nvmlInitWithFlagsFunc nvmlInitWithFlagsFunc = NULL;
static nvmlShutdownFunc nvmlShutdownFunc = NULL;
static nvmlErrorStringFunc nvmlErrorStringFunc = NULL;
static nvmlDeviceGetCountFunc nvmlDeviceGetCountFunc = NULL;
static nvmlDeviceGetHandleByIndexFunc nvmlDeviceGetHandleByIndexFunc = NULL;
static nvmlDeviceGetHandleByUUIDFunc nvmlDeviceGetHandleByUUIDFunc = NULL;
static nvmlDeviceGetMemoryInfoV2Func nvmlDeviceGetMemoryInfoV2Func = NULL;
static nvmlDeviceGetNameFunc nvmlDeviceGetNameFunc = NULL;
static nvmlDeviceGetUUIDFunc nvmlDeviceGetUUIDFunc = NULL;
static nvmlDeviceGetIndexFunc nvmlDeviceGetIndexFunc = NULL;
static nvmlDeviceRegisterEventsFunc nvmlDeviceRegisterEventsFunc = NULL;
static nvmlEventSetCreateFunc nvmlEventSetCreateFunc = NULL;
static nvmlEventSetWaitFunc nvmlEventSetWaitFunc = NULL;
static nvmlEventSetFreeFunc nvmlEventSetFreeFunc = NULL;
static nvmlDeviceGetUtilizationRatesFunc nvmlDeviceGetUtilizationRatesFunc = NULL;
static nvmlDeviceGetComputeRunningProcessesFunc nvmlDeviceGetComputeRunningProcessesFunc = NULL;
static nvmlDeviceGetProcessUtilizationFunc nvmlDeviceGetProcessUtilizationFunc = NULL;
static nvmlDeviceGetMultiGpuBoardFunc nvmlDeviceGetMultiGpuBoardFunc = NULL;
static nvmlDeviceGetTopologyCommonAncestorFunc nvmlDeviceGetTopologyCommonAncestorFunc = NULL;
static nvmlDeviceGetTopologyNearestGpusFunc nvmlDeviceGetTopologyNearestGpusFunc = NULL;
static nvmlSystemGetDriverVersionFunc nvmlSystemGetDriverVersionFunc = NULL;
static nvmlSystemGetCudaDriverVersionFunc nvmlSystemGetCudaDriverVersionFunc = NULL;
static nvmlDeviceGetTemperatureFunc nvmlDeviceGetTemperatureFunc = NULL;
static nvmlDeviceGetPowerUsageFunc nvmlDeviceGetPowerUsageFunc = NULL;

// In order not to depend on libnvidia-ml.so.1, the custom function is implemented as follows:
nvmlReturn_t nvmlInitWrapper(void) {
    return (nvmlInitFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlInitFunc();
}

nvmlReturn_t nvmlInitWithFlagsWrapper(unsigned int flags) {
    return (nvmlInitWithFlagsFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlInitWithFlagsFunc(flags);
}

nvmlReturn_t nvmlShutdownWrapper(void) {
    return (nvmlShutdownFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlShutdownFunc();
}

const char* nvmlErrorStringWrapper(nvmlReturn_t result) {
    if (nvmlErrorStringFunc == NULL) {
        fprintf(stderr, "Could not load function nvmlErrorString from nvml.so\n");
        return "NVML_ERROR_FUNCTION_NOT_FOUND";
    }
    return nvmlErrorStringFunc(result);
}

nvmlReturn_t nvmlSystemGetDriverVersionWrapper(char *version, unsigned int length) {
    return (nvmlSystemGetDriverVersionFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlSystemGetDriverVersionFunc(version, length);
}

nvmlReturn_t nvmlSystemGetCudaDriverVersionWrapper(int *cudaDriverVersion) {
    return (nvmlSystemGetCudaDriverVersionFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlSystemGetCudaDriverVersionFunc(cudaDriverVersion);
}

nvmlReturn_t nvmlDeviceGetTemperatureWrapper(nvmlDevice_t device, nvmlTemperatureSensors_t sensorType, unsigned int *temp) {
    return (nvmlDeviceGetTemperatureFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetTemperatureFunc(device, sensorType, temp);
}

nvmlReturn_t nvmlDeviceGetPowerUsageWrapper(nvmlDevice_t device, unsigned int *power) {
    return (nvmlDeviceGetPowerUsageFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetPowerUsageFunc(device, power);
}

nvmlReturn_t nvmlDeviceGetCountWrapper(unsigned int *deviceCount) {
    return (nvmlDeviceGetCountFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetCountFunc(deviceCount);
}

nvmlReturn_t nvmlDeviceGetHandleByIndexWrapper(unsigned int index, nvmlDevice_t *device) {
    return (nvmlDeviceGetHandleByIndexFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetHandleByIndexFunc(index, device);
}

nvmlReturn_t nvmlDeviceGetHandleByUUIDWrapper(const char *uuid, nvmlDevice_t *device) {
    return (nvmlDeviceGetHandleByUUIDFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetHandleByUUIDFunc(uuid, device);
}

nvmlReturn_t nvmlDeviceGetMemoryInfoV2Wrapper(nvmlDevice_t device, nvmlMemory_v2_t *memory) {
    return (nvmlDeviceGetMemoryInfoV2Func == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetMemoryInfoV2Func(device, memory);
}

nvmlReturn_t nvmlDeviceGetNameWrapper(nvmlDevice_t device, char *name, unsigned int length) {
    return (nvmlDeviceGetNameFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetNameFunc(device, name, length);
}

nvmlReturn_t nvmlDeviceGetUUIDWrapper(nvmlDevice_t device, char *uuid, unsigned int length) {
    return (nvmlDeviceGetUUIDFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetUUIDFunc(device, uuid, length);
}

nvmlReturn_t nvmlDeviceGetIndexWrapper(nvmlDevice_t device, unsigned int *index) {
    return (nvmlDeviceGetIndexFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetIndexFunc(device, index);
}

nvmlReturn_t nvmlDeviceRegisterEventsWrapper(nvmlDevice_t device, unsigned long long eventTypes, nvmlEventSet_t set) {
    return (nvmlDeviceRegisterEventsFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceRegisterEventsFunc(device, eventTypes, set);
}

nvmlReturn_t nvmlEventSetCreateWrapper(nvmlEventSet_t *set) {
    return (nvmlEventSetCreateFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlEventSetCreateFunc(set);
}

nvmlReturn_t nvmlEventSetWaitWrapper(nvmlEventSet_t set, nvmlEventData_t *data, unsigned int timeouts) {
    return (nvmlEventSetWaitFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlEventSetWaitFunc(set, data, timeouts);
}

nvmlReturn_t nvmlEventSetFreeWrapper(nvmlEventSet_t set) {
    return (nvmlEventSetFreeFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlEventSetFreeFunc(set);
}

nvmlReturn_t nvmlDeviceGetUtilizationRatesWrapper(nvmlDevice_t device, nvmlUtilization_t *utilization) {
    return (nvmlDeviceGetUtilizationRatesFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetUtilizationRatesFunc(device, utilization);
}

nvmlReturn_t nvmlDeviceGetComputeRunningProcessesWrapper(nvmlDevice_t device, unsigned int *infoSize, nvmlProcessInfo_v2_t *infos) {
    return (nvmlDeviceGetComputeRunningProcessesFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetComputeRunningProcessesFunc(device, infoSize, infos);
}

nvmlReturn_t nvmlDeviceGetProcessUtilizationWrapper(nvmlDevice_t device, nvmlProcessUtilizationSample_t *samples, unsigned int sampleSize, unsigned long long timestamp) {
    return (nvmlDeviceGetProcessUtilizationFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetProcessUtilizationFunc(device, samples, sampleSize, timestamp);
}

nvmlReturn_t nvmlDeviceGetMultiGpuBoardWrapper(nvmlDevice_t device, unsigned int *multiGpuBoard) {
    return (nvmlDeviceGetMultiGpuBoardFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetMultiGpuBoardFunc(device, multiGpuBoard);
}

nvmlReturn_t nvmlDeviceGetTopologyCommonAncestorWrapper(nvmlDevice_t device1, nvmlDevice_t device2, nvmlGpuTopologyLevel_t *pathInfo) {
    return (nvmlDeviceGetTopologyCommonAncestorFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetTopologyCommonAncestorFunc(device1, device2, pathInfo);
}

nvmlReturn_t nvmlDeviceGetTopologyNearestGpusWrapper(nvmlDevice_t device, nvmlGpuTopologyLevel_t level, unsigned int *count, nvmlDevice_t *deviceArray) {
    return (nvmlDeviceGetTopologyNearestGpusFunc == NULL) ? NVML_ERROR_FUNCTION_NOT_FOUND : nvmlDeviceGetTopologyNearestGpusFunc(device, level, count, deviceArray);
}

// Helper function to load a symbol and handle errors.
static int loadSymbol(const char *symbolName, void **symbolPtr) {
    *symbolPtr = dlsym(nvmlHandle, symbolName);
    if (*symbolPtr == NULL) {
        fprintf(stderr, "Failed to load symbol %s\n", symbolName);
        return 0;
    }
    return 1;
}

// Loads the "libnvidia-ml.so.1" shared library and all required symbols.
static nvmlReturn_t loadLibrary(void) {
    nvmlHandle = dlopen(nvmlLibraryName, RTLD_LAZY);
    if (nvmlHandle == NULL) {
        fprintf(stderr, "Failed to load symbol libnvidia-ml.so.1: %s\n", dlerror());
        return NVML_ERROR_LIBRARY_NOT_FOUND;
    }

    loadSymbol("nvmlInit_v2", (void**)&nvmlInitFunc);
    loadSymbol("nvmlInitWithFlags_v2", (void**)&nvmlInitWithFlagsFunc);
    loadSymbol("nvmlShutdown", (void**)&nvmlShutdownFunc);
    loadSymbol("nvmlErrorString", (void**)&nvmlErrorStringFunc);
    loadSymbol("nvmlDeviceGetCount_v2", (void**)&nvmlDeviceGetCountFunc);
    loadSymbol("nvmlDeviceGetHandleByIndex_v2", (void**)&nvmlDeviceGetHandleByIndexFunc);
    loadSymbol("nvmlDeviceGetHandleByUUID_v2", (void**)&nvmlDeviceGetHandleByUUIDFunc);
    loadSymbol("nvmlDeviceGetMemoryInfo_v2", (void**)&nvmlDeviceGetMemoryInfoV2Func);
    loadSymbol("nvmlDeviceGetName_v2", (void**)&nvmlDeviceGetNameFunc);
    loadSymbol("nvmlDeviceGetUUID_v2", (void**)&nvmlDeviceGetUUIDFunc);
    loadSymbol("nvmlDeviceGetIndex", (void**)&nvmlDeviceGetIndexFunc);
    loadSymbol("nvmlDeviceRegisterEvents", (void**)&nvmlDeviceRegisterEventsFunc);
    loadSymbol("nvmlEventSetCreate", (void**)&nvmlEventSetCreateFunc);
    loadSymbol("nvmlEventSetWait_v2", (void**)&nvmlEventSetWaitFunc);
    loadSymbol("nvmlEventSetFree", (void**)&nvmlEventSetFreeFunc);
    loadSymbol("nvmlDeviceGetUtilizationRates", (void**)&nvmlDeviceGetUtilizationRatesFunc);
    loadSymbol("nvmlDeviceGetComputeRunningProcesses_v2", (void**)&nvmlDeviceGetComputeRunningProcessesFunc);
    loadSymbol("nvmlDeviceGetProcessUtilization", (void**)&nvmlDeviceGetProcessUtilizationFunc);
    loadSymbol("nvmlDeviceGetMultiGpuBoard", (void**)&nvmlDeviceGetMultiGpuBoardFunc);
    loadSymbol("nvmlDeviceGetTopologyCommonAncestor", (void**)&nvmlDeviceGetTopologyCommonAncestorFunc);
    loadSymbol("nvmlDeviceGetTopologyNearestGpus", (void**)&nvmlDeviceGetTopologyNearestGpusFunc);
    loadSymbol("nvmlSystemGetDriverVersion", (void**)&nvmlSystemGetDriverVersionFunc);
    loadSymbol("nvmlSystemGetCudaDriverVersion", (void**)&nvmlSystemGetCudaDriverVersionFunc);
    loadSymbol("nvmlDeviceGetTemperature", (void**)&nvmlDeviceGetTemperatureFunc);
    loadSymbol("nvmlDeviceGetPowerUsage", (void**)&nvmlDeviceGetPowerUsageFunc);

    fprintf(stdout, "Load libnvidia-ml.so.1 success!\n");
    return NVML_SUCCESS;
}

// Shut down NVML and decrements the reference count on the dynamically loaded
// "libnvidia-ml.so.1" library.
// Call this once no NVML is no longer being used.
static nvmlReturn_t unloadLibrary(void) {
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

var errNvmlLoaded = errors.New("Could not load NVML library succeed")

// adapter for go-own api
type cnvmlDevice struct {
	Handle C.nvmlDevice_t
}

type cnvmlEventSet struct {
	Handle C.nvmlEventSet_t
}

// Convert the nvmlReturn_t return value to a more readable error message
func errorString(ret C.nvmlReturn_t) error {
	if ret == C.NVML_SUCCESS {
		return nil
	}
	if C.nvmlHandle == nil {
		log.Printf("Can't load function nvmlErrorString from nvml.so")
		return errNvmlLoaded
	}
	err := C.GoString(C.nvmlErrorStringWrapper(ret))
	return fmt.Errorf("nvml: %v", err)
}

var cgoAllocUnknown = new(struct{}) // Used to clear 'assignment mismatch' alarms. Refer to NVIDIA official code.

// The way strings are represented inside the Go Language
type StringHeader struct {
	Data unsafe.Pointer
	Len  uint
}

// unpackCString represents the data from Go string as *C.char and avoids copying.
func unpackCString(str string) (*C.char, *struct{}) {
	sh := (*StringHeader)(unsafe.Pointer(&str))
	return (*C.char)(sh.Data), cgoAllocUnknown
}

func loadNvmlSo() error {
	err := errorString(C.loadLibrary())
	if err != nil {
		log.Printf("loadNvmlSo failed: %v", err)
		return err
	}
	return nil
}

func unloadNvmlSo() error {
	return errorString(C.unloadLibrary())
}

func nvmlInitWrapper() NvmlRetType {
	return NvmlRetType(C.nvmlInitWrapper())
}

func nvmlInitWithFlagsWrapper(flags uint32) NvmlRetType {
	return NvmlRetType(C.nvmlInitWithFlagsWrapper(C.uint(flags)))
}

func nvmlShutdownWrapper() NvmlRetType {
	return NvmlRetType(C.nvmlShutdownWrapper())
}

func nvmlDeviceGetCountWrapper(DeviceCount *uint32) NvmlRetType {
	return NvmlRetType(C.nvmlDeviceGetCountWrapper((*C.uint)(unsafe.Pointer(DeviceCount))))
}

func nvmlDeviceGetCountWrapper(DeviceCount *uint32) NvmlRetType {
    cDeviceCount, _ := (*C.uint)(unsafe.Pointer(DeviceCount)), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetCount(cDeviceCount))
}

func nvmlSystemGetDriverVersionWrapper(Version *byte, Length uint32) NvmlRetType {
    cVersion, _ := (*C.char)(unsafe.Pointer(Version)), cgoAllocUnknown
    cLength, _ := C.uint(Length), cgoAllocUnknown
    return NvmlRetType(C.nvmlSystemGetDriverVersion(cVersion, cLength))
}

func nvmlSystemGetCudaDriverVersionWrapper(CudaDriverVersion *int32) NvmlRetType {
    cCudaDriverVersion, _ := (*C.int)(unsafe.Pointer(CudaDriverVersion)), cgoAllocUnknown
    return NvmlRetType(C.nvmlSystemGetCudaDriverVersion(cCudaDriverVersion))
}

func nvmlDeviceGetHandleByIndexWrapper(Index uint32, NvmlDevice *NvmlDevice) NvmlRetType {
    cIndex, _ := C.uint(Index), cgoAllocUnknown
    cNvmlDevice, _ := (*C.nvmlDevice_t)(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetHandleByIndex(cIndex, cNvmlDevice))
}

func nvmlDeviceGetHandleByUUIDWrapper(Uuid string, NvmlDevice *NvmlDevice) NvmlRetType {
    cUuid, _ := unpackCString(Uuid)
    cNvmlDevice, _ := (*C.nvmlDevice_t)(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetHandleByUUID(cUuid, cNvmlDevice))
}

func nvmlDeviceGetMemoryInfoV2Wrapper(NvmlDevice NvmlDevice, Memory *MemoryV2) NvmlRetType {
    cNvmlDevice, _ := C.nvmlDevice_t(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    cMemory, _ := (*C.nvmlMemory_v2_t)(unsafe.Pointer(Memory)), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetMemoryInfoV2(cNvmlDevice, cMemory))
}

func nvmlDeviceGetNameWrapper(NvmlDevice NvmlDevice, Name *byte, Length uint32) NvmlRetType {
    cNvmlDevice, _ := C.nvmlDevice_t(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    cName, _ := (*C.char)(unsafe.Pointer(Name)), cgoAllocUnknown
    cLength, _ := C.uint(Length), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetName(cNvmlDevice, cName, cLength))
}

func nvmlDeviceGetUUIDWrapper(NvmlDevice NvmlDevice, Uuid *byte, Length uint32) NvmlRetType {
    cNvmlDevice, _ := C.nvmlDevice_t(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    cUuid, _ := (*C.char)(unsafe.Pointer(Uuid)), cgoAllocUnknown
    cLength, _ := C.uint(Length), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetUUID(cNvmlDevice, cUuid, cLength))
}

func nvmlDeviceGetIndexWrapper(NvmlDevice NvmlDevice, Index *uint32) NvmlRetType {
    cNvmlDevice, _ := C.nvmlDevice_t(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    cIndex, _ := (*C.uint)(unsafe.Pointer(Index)), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetIndex(cNvmlDevice, cIndex))
}

func nvmlDeviceRegisterEventsWrapper(NvmlDevice NvmlDevice, EventTypes uint64, Set NvmlEventSet) NvmlRetType {
    cNvmlDevice, _ := C.nvmlDevice_t(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    cEventTypes, _ := C.ulonglong(EventTypes), cgoAllocUnknown
    cSet, _ := C.nvmlEventSet_t(unsafe.Pointer(Set)), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceRegisterEvents(cNvmlDevice, cEventTypes, cSet))
}

func nvmlEventSetCreateWrapper(Set *NvmlEventSet) NvmlRetType {
    cSet, _ := (*C.nvmlEventSet_t)(unsafe.Pointer(Set)), cgoAllocUnknown
    return NvmlRetType(C.nvmlEventSetCreate(cSet))
}

func nvmlEventSetWaitWrapper(Set NvmlEventSet, Data *NvmlEventData, Timeouts uint32) NvmlRetType {
    cSet, _ := C.nvmlEventSet_t(unsafe.Pointer(Set)), cgoAllocUnknown
    cData, _ := (*C.nvmlEventData_t)(unsafe.Pointer(Data)), cgoAllocUnknown
    cTimeouts, _ := C.uint(Timeouts), cgoAllocUnknown
    return NvmlRetType(C.nvmlEventSetWait(cSet, cData, cTimeouts))
}

func nvmlEventSetFreeWrapper(Set NvmlEventSet) NvmlRetType {
    cSet, _ := C.nvmlEventSet_t(unsafe.Pointer(Set)), cgoAllocUnknown
    return NvmlRetType(C.nvmlEventSetFree(cSet))
}

func nvmlDeviceGetUtilizationRatesWrapper(NvmlDevice NvmlDevice, Utilization *Utilization) NvmlRetType {
    cNvmlDevice, _ := C.nvmlDevice_t(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    cUtilization, _ := (*C.nvmlUtilization_t)(unsafe.Pointer(Utilization)), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetUtilizationRates(cNvmlDevice, cUtilization))
}

func nvmlDeviceGetComputeRunningProcessesWrapper(NvmlDevice NvmlDevice, InfoCount *uint32, Infos []ProcessInfoV2) NvmlRetType {
    cNvmlDevice, _ := C.nvmlDevice_t(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    cInfoCount, _ := (*C.uint)(unsafe.Pointer(InfoCount)), cgoAllocUnknown
    cInfos, _ := (*C.nvmlProcessInfo_v2_t)(unsafe.Pointer(&Infos[0])), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetComputeRunningProcesses(cNvmlDevice, cInfoCount, cInfos))
}

func nvmlDeviceGetProcessUtilizationWrapper(NvmlDevice NvmlDevice, Utilization []ProcessUtilizationSample, SampleCount uint32, LastGpuTimestamp uint64) NvmlRetType {
    cNvmlDevice, _ := C.nvmlDevice_t(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    cUtilization, _ := (*C.nvmlProcessUtilizationSample_t)(unsafe.Pointer(&Utilization[0])), cgoAllocUnknown
    cSampleCount, _ := C.uint(SampleCount), cgoAllocUnknown
    cLastGpuTimestamp, _ := C.ulonglong(LastGpuTimestamp), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetProcessUtilization(cNvmlDevice, cUtilization, cSampleCount, cLastGpuTimestamp))
}

func nvmlDeviceGetMultiGpuBoardWrapper(NvmlDevice NvmlDevice, MultiGpuBool *uint32) NvmlRetType {
    cNvmlDevice, _ := C.nvmlDevice_t(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    cMultiGpuBool, _ := (*C.uint)(unsafe.Pointer(MultiGpuBool)), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetMultiGpuBoard(cNvmlDevice, cMultiGpuBool))
}

func nvmlDeviceGetTopologyCommonAncestorWrapper(Device1 NvmlDevice, Device2 NvmlDevice, PathInfo *GpuTopologyLevel) NvmlRetType {
    cDevice1, _ := C.nvmlDevice_t(unsafe.Pointer(Device1)), cgoAllocUnknown
    cDevice2, _ := C.nvmlDevice_t(unsafe.Pointer(Device2)), cgoAllocUnknown
    cPathInfo, _ := (*C.nvmlGpuTopologyLevel_t)(unsafe.Pointer(PathInfo)), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetTopologyCommonAncestor(cDevice1, cDevice2, cPathInfo))
}

func nvmlDeviceGetTopologyNearestGpusWrapper(NvmlDevice NvmlDevice, Level GpuTopologyLevel, Count *uint32, DeviceArray []NvmlDevice) NvmlRetType {
    cNvmlDevice, _ := C.nvmlDevice_t(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    cLevel, _ := C.nvmlGpuTopologyLevel_t(Level), cgoAllocUnknown
    cCount, _ := (*C.uint)(unsafe.Pointer(Count)), cgoAllocUnknown
    cDeviceArray, _ := (*C.nvmlDevice_t)(unsafe.Pointer(&DeviceArray[0])), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetTopologyNearestGpus(cNvmlDevice, cLevel, cCount, cDeviceArray))
}

func nvmlDeviceGetTemperatureWrapper(NvmlDevice NvmlDevice, NvmlTemperatureGpu NvmlTemperatureSensors, Temp *uint32) NvmlRetType {
    cNvmlDevice, _ := C.nvmlDevice_t(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    cNvmlTemperatureGpu, _ := C.nvmlTemperatureSensors_t(NvmlTemperatureGpu), cgoAllocUnknown
    cTemp, _ := (*C.uint)(unsafe.Pointer(Temp)), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetTemperature(cNvmlDevice, cNvmlTemperatureGpu, cTemp))
}

func nvmlDeviceGetPowerUsageWrapper(NvmlDevice NvmlDevice, Power *uint32) NvmlRetType {
    cNvmlDevice, _ := C.nvmlDevice_t(unsafe.Pointer(NvmlDevice)), cgoAllocUnknown
    cPower, _ := (*C.uint)(unsafe.Pointer(Power)), cgoAllocUnknown
    return NvmlRetType(C.nvmlDeviceGetPowerUsage(cNvmlDevice, cPower))
}