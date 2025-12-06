/*
Copyright (c) Huawei Technologies Co., Ltd. 2024-2025. All rights reserved.
*/

// Package service implements service of getting pids
package service

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gommonkey/v2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"huawei.com/xpu-device-plugin/pkg/plugin/types"
	"huawei.com/xpu-device-plugin/pkg/plugin/util"
)

const (
	pid1       = 27761708
	pid2       = 25934978
	processMem = 68
)

func isSliceEqual(a, b interface{}) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	sliceA := reflect.ValueOf(a)
	sliceB := reflect.ValueOf(b)
	if sliceA.Len() != sliceB.Len() {
		return false
	}
	for i := 0; i < sliceA.Len(); i++ {
		if !reflect.DeepEqual(sliceA.Index(i).Interface(), sliceB.Index(i).Interface()) {
			return false
		}
	}
	return true
}

func TestReadProcsFile(t *testing.T) {
	realHostPids, err := readProcsFile("./test_cgroup.procs")
	expectedHostPids := []int{27761708, 25934978}
	if err == nil && isSliceEqual(realHostPids, expectedHostPids) {
		t.Log("test succeed")
	} else {
		t.Error("test failed")
	}
}

func TestReadStatusFile(t *testing.T) {
	realPids, err := readStatusFile("./test.status")
	expectedPids := fmt.Sprintf("%d %d", 27761708, 1)
	if err == nil && realPids == expectedPids {
		t.Log("test succeed")
	} else {
		t.Error("test failed")
	}
}

func TestParseDockerCgroupPath(t *testing.T) {
	cgroupPath := "/kubepods/besteffort/pod219299e_6f3c_475e_98c7_fa2d7bc53676" +
		"/533d32462f5c2573970e217e018d653348a852c29217ec29a1188"
	realPodId, realContainerId, err := parseDockerCgroupPath(cgroupPath)
	expectedPodId := "219299e-6f3c-475e-98c7-fa2d7bc53676"
	expectedContainerId := "533d32462f5c2573970e217e018d653348a852c29217ec29a1188"
	if err == nil && realPodId == expectedPodId && realContainerId == expectedContainerId {
		t.Log("test succeed")
	} else {
		t.Error("test failed")
	}
}

func TestParseCgroupPath(t *testing.T) {
	cgroupPath := "kubepods/besteffort/pod4bd40d_4fd5_4691_95ab_dbce6e9fa1db.slice" +
		"/cri-containerd-4e331f94fb115736682f59d310e517c7f227f1d3d581724362885ec2544d5.scope"
	realPodId, realContainerId, err := parseCgroupPath(cgroupPath)
	expectedPodId := "4bd40d-4fd5-4691-95ab-dbce6e9fa1db"
	expectedContainerId := "4e331f94fb115736682f59d310e517c7f227f1d3d581724362885ec2544d5"
	if err == nil && realPodId == expectedPodId && realContainerId == expectedContainerId {
		t.Log("test succeed")
	} else {
		t.Error("test failed")
	}
}

func TestGetContainerName(t *testing.T) {
	cgroupPath := "/kubepods/besteffort/pod4bd40d_4fd5_4691_95ab_dbce6e9fa1db.slice" +
		"/cri-containerd-4e331f94fb115736682f59d310e517c7f227f1d3d581724362885ec2544d5.scope"
	patch := gommonkey.ApplyFunc(util.ListPods, func(opts metav1.ListOptions) (*v1.PodList, error) {
		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{UID: "4bd40d-4fd5-4691-95ab-dbce6e9fa1db"},
			Status:     v1.PodStatus{Phase: v1.PodRunning},
		}
		cs := v1.ContainerStatus{
			ContainerID: "containerd://4e331f94fb115736682f59d310e517c7f227f1d3d581724362885ec2544d5",
			Ready:       true,
			Name:        "container0",
		}
		pod.Status.ContainerStatuses = []v1.ContainerStatus{cs}
		return &v1.PodList{Items: []v1.Pod{pod}}, nil
	})
	defer patch.Reset()

	realPodId, realContainerName, err := getContainerName(cgroupPath)
	expectedPodId := "4bd40d-4fd5-4691-95ab-dbce6e9fa1db"
	expectedContainerName := "container0"
	if err == nil && realPodId == expectedPodId && realContainerName == expectedContainerName {
		t.Log("test succeed")
	} else {
		t.Error("test failed")
	}
}

func fileCheckSum(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func TestWritePidsConfig(t *testing.T) {
	pidMaps := fmt.Sprintf("%-15s %-15s", "27761708", "1", "2934978", "7")
	err := writePidsConfig("./real.pids.config", pidMaps)
	if err == nil && fileCheckSum("./real.pids.config") == fileCheckSum("./expected.pids.config") {
		t.Log("test succeed")
	} else {
		t.Error("test failed")
	}
}

func TestReadPidsConfig(t *testing.T) {
	pids, err := readPidsConfig("./expected.pids.config")
	expectedHostPids := []uint32{27761708, 2934978}
	if err == nil && isSliceEqual(pids, expectedHostPids) {
		t.Log("test succeed")
	} else {
		t.Error("test failed")
	}
}

func TestCleanDestroyedPodDir(t *testing.T) {
	patchGetPodDirNames := gommonkey.ApplyFunc(getPodDirNames, func() ([]string, error) {
		return []string{"d1640c4c-6389-bd7d-7b164f7c7c4", "1abab0c4-4fd5-4691-95ab-dbce6e9fa1db"}, nil
	})
	defer patchGetPodDirNames.Reset()

	patchListPods := gommonkey.ApplyFunc(util.ListPods, func(opts metav1.ListOptions) (*v1.PodList, error) {
		pod := v1.Pod{}
		pod.UID = "1abab0c4-4fd5-4691-95ab-dbce6e9fa1db"
		podList := &v1.PodList{Items: []v1.Pod{pod}}
		return podList, nil
	})
	defer patchListPods.Reset()

	err := cleanDestroyedPodDir()
	if err == nil {
		t.Log("test succeed")
	} else {
		t.Error("test failed")
	}
}

func makeXpuDeviceAndVgpuDevices() (map[string]*types.XpuDevice, types.VgpuDevices) {
	xpuDevices := make(map[string]*types.XpuDevice)
	device := types.XpuDevice{
		Index:  6,
		ID:     "id0",
		Type:   "NVIDIA",
		Count:  20,
		Memory: types.Memory{Total: 200},
		Health: true,
		VgpuDeviceList: types.VgpuDevices{},
	}
	xpuDevices["id0"] = &device

	vgpuDevices := types.VgpuDevices{}
	dev0 := types.VgpuDevice{
		Id:              "id0-0",
		Gpuid:           "id0",
		PodUID:          "pod0",
		ContainerID:     "container0",
		ContainerName:   "container0",
		VgpuMemoryLimit: 1 * 1024,
		VgpuCoreLimit:   5,
	}
	dev1 := types.VgpuDevice{
		Id:              "id0-1",
		Gpuid:           "id0",
		PodUID:          "pod1",
		ContainerID:     "container1",
		ContainerName:   "container1",
		VgpuMemoryLimit: 1 * 1024,
		VgpuCoreLimit:   5,
	}
	vgpuDevices = append(vgpuDevices, dev0)
	vgpuDevices = append(vgpuDevices, dev1)
	return xpuDevices, vgpuDevices
}

func TestSetVgpuDevices(t *testing.T) {
	xpuDevices, vgpuDevices := makeXpuDeviceAndVgpuDevices()
	uidToProcessMap := make(map[string]map[uint32]types.ProcessUsage)
	processMap := make(map[uint32]types.ProcessUsage)
	p0 := types.ProcessUsage{ProcessMem: processMem * 1024 * 1024, ProcessCoreUtilization: 15}
	p1 := types.ProcessUsage{ProcessMem: processMem * 1024 * 1024, ProcessCoreUtilization: 20}
	processMap[pid1] = p0
	processMap[pid2] = p1
	uidToProcessMap["id0"] = processMap

	var pids []uint32
	pids1 := []uint32{pid1}
	pSet["pod0/container0"] = pids1
	pids2 := []uint32{pid2}
	pSet["pod1/container1"] = pids2

	expectedVgpuDevices := types.VgpuDevices{}
	expectedDev0 := types.VgpuDevice{
		Id:                   "id0-0",
		Gpuid:                "id0",
		PodUID:               "pod0",
		ContainerName:        "container0",
		VgpuMemoryUsed:       48,
		VgpuCoreUtilization:  15,
		VgpuMemoryLimit:      1 * 1024,
		VgpuCoreLimit:        5,
	}
	expectedDev1 := types.VgpuDevice{
		Id:                   "id0-1",
		Gpuid:                "id0",
		PodUID:               "pod1",
		ContainerName:        "container1",
		VgpuMemoryUsed:       48,
		VgpuCoreUtilization:  20,
		VgpuMemoryLimit:      1 * 1024,
		VgpuCoreLimit:        5,
	}
	expectedVgpuDevices = append(expectedVgpuDevices, expectedDev0)
	expectedVgpuDevices = append(expectedVgpuDevices, expectedDev1)
	xpuDevices = setVgpuDevices(xpuDevices, vgpuDevices, uidToProcessMap, pSet)
	if isSliceEqual(xpuDevices["id0"].VgpuDeviceList, expectedVgpuDevices) &&
		xpuDevices["id0"].MemoryUtilization == 100 &&
		xpuDevices["id0"].MemoryUtilization == 50 {
		t.Log("test succeed")
	} else {
		t.Error("test failed")
	}
}