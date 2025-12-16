/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2025. All rights reserved.
 */

// Package util implements util function for device plugin
package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"huawei.com/xpu-device-plugin/pkg/lock"
	"huawei.com/xpu-device-plugin/pkg/plugin/config"
	"huawei.com/xpu-device-plugin/pkg/plugin/types"
	"huawei.com/xpu-device-plugin/pkg/plugin/xpu"
)

const (
	NodeLength                 = 7
	// PodAnnotationMaxLength annotation max data length 2MB
	PodAnnotationMaxLength     = 1024 * 1024
	BaseDec                    = 10
	BitsSize                   = 64
	// BitSize base size
)

func init() {
	lock.NewClient()
}

// GetNode get k8s node object according to node name
func GetNode(nodeName string) (*v1.Node, error) {
	return lock.GetClient().CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
}

// ListPods list k8s pods according to list options
func ListPods(opts metav1.ListOptions) (*v1.PodList, error) {
	return lock.GetClient().CoreV1().Pods("").List(context.Background(), opts)
}

// GetPendingPod get k8s pod object according to node name and types.DeviceBindAllocating status
func GetPendingPod(nodeName string) (*v1.Pod, error) {
	podList, err := ListPods(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var oldestPod *v1.Pod
	oldestBindTime := uint64(math.MaxUint64)

	for _, p := range podList.Items {
		bindTime, ok := getBindTime(p)
		if !ok {
			continue
		}

		phase, ok := p.Annotations[types.DeviceBindPhase]
		if !ok {
			continue
		} else if strings.Compare(phase, types.DeviceBindAllocating) != 0 {
			continue
		}

		n, ok := p.Annotations[xpu.AssignedNode]
		if !ok {
			continue
		} else if strings.Compare(n, nodeName) != 0 {
			continue
		}

		if oldestBindTime > bindTime {
			oldestBindTime = bindTime
			oldestPod = &p
		}
	}
	return oldestPod, nil
}

func getBindTime(pod v1.Pod) (uint64, bool) {
	assumeTimeStr, ok := pod.Annotations[types.DeviceBindTime]
	if !ok {
		return math.MaxUint64, false
	}

	if len(assumeTimeStr) > PodAnnotationMaxLength {
		log.Warningf("timestamp annotation invalid, pod Name: %s", pod.Name)
		return math.MaxUint64, false
	}

	bindTime, err := strconv.ParseUint(assumeTimeStr, BaseDec, BitsSize)
	if err != nil {
		log.Errorln("parse timestamp failed, %v", err)
		return math.MaxUint64, false
	}
	return bindTime, true
}

// EncodeNodeDevices encode a node's xpus info to string
func EncodeNodeDevices(devList []types.DeviceInfo) string {
	var encodedNodeDevices strings.Builder
	for _, val := range devList {
		encodedNodeDevices.Write([]byte(strconv.Itoa(int(val.Index))))
		encodedNodeDevices.Write([]byte(","))
		encodedNodeDevices.Write([]byte(val.Id))
		encodedNodeDevices.Write([]byte(","))
		encodedNodeDevices.Write([]byte(strconv.Itoa(int(val.Count))))
		encodedNodeDevices.Write([]byte(","))
		encodedNodeDevices.Write([]byte(strconv.Itoa(int(val.Devmem))))
		encodedNodeDevices.Write([]byte(","))
		encodedNodeDevices.Write([]byte(val.Type))
		encodedNodeDevices.Write([]byte(","))
		encodedNodeDevices.Write([]byte(fmt.Sprintf("%v", val.Health)))
		encodedNodeDevices.Write([]byte(","))
		encodedNodeDevices.Write([]byte(strconv.Itoa(int(val.Numa))))
		encodedNodeDevices.Write([]byte(";"))
	}
	log.Infoln("Encoded node Devices: ", encodedNodeDevices.String())
	return encodedNodeDevices.String()
}

// EncodeContainerDevices encode xpu resource request of a container to string
func EncodeContainerDevices(cd types.ContainerDevices) string {
	var encodedContainerDevices strings.Builder
	for _, val := range cd {
		encodedContainerDevices.Write([]byte(strconv.Itoa(int(val.Index))))
		encodedContainerDevices.Write([]byte(","))
		encodedContainerDevices.Write([]byte(val.UUID))
		encodedContainerDevices.Write([]byte(","))
		encodedContainerDevices.Write([]byte(val.Type))
		encodedContainerDevices.Write([]byte(","))
		encodedContainerDevices.Write([]byte(strconv.Itoa(int(val.Usedmem))))
		encodedContainerDevices.Write([]byte(","))
		encodedContainerDevices.Write([]byte(strconv.Itoa(int(val.UsedCores))))
		encodedContainerDevices.Write([]byte(","))
		encodedContainerDevices.Write([]byte(strconv.Itoa(int(val.Vid))))
		encodedContainerDevices.Write([]byte(";"))
	}
	log.Infoln("Encoded container Devices: ", encodedContainerDevices.String())
	return encodedContainerDevices.String()
}

// EncodePodDevices encode xpu resource request of a pod to string
func EncodePodDevices(pd types.PodDevices) string {
	if len(pd) == 0 {
		return ""
	}
	var es []string
	for _, cd := range pd {
		es = append(es, EncodeContainerDevices(cd))
	}
	return strings.Join(es, ",")
}

// GetXpuDevice get XPUdevice info
func GetXpuDevice(str string) map[string]map[string]*types.XPUDevice {
	deviceMap := make(map[string]map[string]*types.XPUDevice)
	driverVersion, FrameworkVersion, err := xpu.GetVersionInfo()
	if err != nil {
		log.infof("getVersionInfo error %v", err)
	}
	for _, deviceInfo := range deviceMap {
		deviceInfo.NodeIp = ip
		deviceInfo.NodeName = config.NodeName
		deviceInfo.DriverVersion = driverVersion
		deviceInfo.FrameworkVersion = FrameworkVersion
	}
	return deviceMap
}

// DecodeNodeDevices decode the node device from string
func DecodeNodeDevices(str string) map[string]map[string]*types.XPUDevice {
	deviceMap := make(map[string]map[string]*types.XPUDevice)
	if !strings.Contains(str, ";") {
		log.Errorln("decode node device failed, wrong annos: %s", str)
		return deviceMap
	}
	tmp := strings.Split(str, ";")
	for _, val := range tmp {
		if !strings.Contains(val, ",") {
			continue
		}
		items := strings.Split(val, ",")
		if len(items) != 7 {
			log.Warningf("device string is wrong, device: %s", items)
			continue
		}
		index, err := strconv.Atoi(items[0])
		if err != nil {
			continue
		}
		count, err := strconv.Atoi(items[2])
		if err != nil {
			continue
		}
		devmem, err := strconv.Atoi(items[3])
		if err != nil {
			continue
		}
		health, err := strconv.ParseBool(items[5])
		if err != nil {
			continue
		}
		i := types.XPUDevice{
			Index:          int32(index),
			Id:             items[1],
			Type:           items[4],
			Count:          uint32(count),
			MemoryTotal:    uint64(devmem),
			Health:         health,
			VxpuDeviceList: types.VxpuDevices{},
		}
		deviceMap[items[1]] = &i
	}
	return deviceMap
}

// DecodeContainerDevices decode xpu resource request of a container from string
func DecodeContainerDevices(str string) types.ContainerDevices {
	if len(str) == 0 {
		return types.ContainerDevices{}
	}
	cd := types.ContainerDevices{}
	content := strings.Split(str, ";")
	for _, val := range content {
		if strings.Contains(val, ",") == false {
			continue
		}
		fields := strings.Split(val, ",")
		tmpdev := reflect.TypeOf(tmpdev).NumField()
		if len(fields) != tmpdev {
			log.Fatalln("DecodeContainerDevices invalid parameter:", str)
			return types.ContainerDevices{}
		}
		index, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatalln("DecodeContainerDevices invalid parameter:", str)
			return types.ContainerDevices{}
		}
		tmpdev.Index = int32(index)
		tmpdev.UUID = fields[1]
		tmpdev.Type = fields[2]
		devcores, err := strconv.Atoi(fields[3])
		if err != nil {
			log.Fatalln("DecodeContainerDevices invalid parameter:", str)
			return types.ContainerDevices{}
		}
		tmpdev.Usedmem = int32(index)
		devcores, err := strconv.Atoi(fields[4])
		if err != nil {
			log.Fatalln("DecodeContainerDevices invalid parameter:", str)
			return types.ContainerDevices{}
		}
		tmpdev.Usedcores = int32(devcores)
		vid, err := strconv.Atoi(fields[5])
		if err != nil {
			log.Fatalln("DecodeContainerDevices invalid parameter:", str)
			return types.ContainerDevices{}
		}
		tmpdev.Vid = int32(vid)
		contdev = append(contdev, tmpdev)
	}
	return contdev
}

// DecodePodDevices decode xpu resource request of a pod from string
func DecodePodDevices(str string) types.PodDevices {
	if len(str) == 0 {
		return types.PodDevices{}
	}
	var pd types.PodDevices
	for _, s := range strings.Split(str, ",") {
		cd := DecodeContainerDevices(s)
		pd = append(pd, cd)
	}
	return pd
}

func getContainerIdxByVgpuIdx(p *v1.Pod, vgpuIdx int) int {
	foundVgpuIdx := -1
	for i, container := range p.Spec.Containers {
		if ok := container.Resources.Limits[xpu.VgpuNumber]; ok {
			foundVgpuIdx++
			if foundVgpuIdx == vgpuIdx {
				return i
			}
			continue
		}
		if ok := container.Resources.Limits[xpu.VgpuMemory]; ok {
			foundVgpuIdx++
			if foundVgpuIdx == vgpuIdx {
				return i
			}
			continue
		}
	}
	return -1
}

// Get xgpu limit info of the container
func getVgpuLimit(resourcelist v1.ResourceList) (int64, int64, int64) {
	var number int64 = 0
	var core int64 = 0
	var mem int64 = 0
	if vgpuNumber, ok := resourcelist[xpu.VgpuNumber]; ok {
		number = vgpuNumber.Value()
	}
	if vgpuCore, ok := resourcelist[xpu.VgpuCore]; ok {
		core = vgpuCore.Value()
	}
	if vgpuMem, ok := resourcelist[xpu.VgpuMemory]; ok {
		mem = vgpuMem.Value()
	}
	return number, core, mem
}

// GetNextDeviceRequest get next xpu resource request of container in a pod
// reference code: https://gitee.com/openeuler/kubernetes/blob/master/pkg/scheduler/app/plugins/deviceplugin/gpu/util.go
func GetNextDeviceRequest(devType string, p v1.Pod) (v1.Container, types.ContainerDevices, error) {
	podDevices := DecodePodDevices(p.Annotations[xpu.AssignedDevicesToAllocate])
	res := types.ContainerDevices{}
	var found bool
	for _, val := range podDevices {
		for _, dev := range val {
			if strings.Compare(dev.Type, devType) == 0 {
				res = append(res, dev)
				found = true
			}
		}
		if found {
			break
		}
	}
	if !found {
		return v1.Container{}, res, errors.New("device request not found")
	}
	vgpuIdx := 0
	for _, dev := range res {
		if dev.Type == devType {
			vgpuIdx++
		}
	}
	idx := getContainerIdxByVgpuIdx(&p, vgpuIdx)
	if idx == -1 {
		log.Errorln("get container idx by vgpuIdx failed, vgpuIdx: %v", vgpuIdx)
		return v1.Container{}, res, nil
	}
	return p.Spec.Containers[idx], res, nil
}

// EraseNextDeviceTypeFromAnnotation erase next xpu resource request of container in a pod's annotation
func EraseNextDeviceTypeFromAnnotation(devType string, p v1.Pod) error {
	annotations := p.Annotations
	podDevices := DecodePodDevices(annotations[xpu.AssignedDevicesToAllocate])
	res := types.PodDevices{}
	var found bool
	for _, val := range podDevices {
		if found {
			res = append(res, val)
			continue
		}
		tmp := types.ContainerDevices{}
		for _, dev := range val {
			if strings.Compare(dev.Type, devType) == 0 {
				found = true
			} else {
				tmp = append(tmp, dev)
			}
		}
		if found {
			res = append(res, tmp)
		} else {
			res = append(res, val)
		}
	}
	newAnnos := make(map[string]string)
	newAnnos[xpu.AssignedDevicesToAllocate] = EncodePodDevices(res)
	log.Infoln("After erase res is :", EncodePodDevices(res))
	return PatchPodAnnotations(p, newAnnos)
}

// PodAllocationSuccess try to patch annotation of a pod to indicate allocation success
func PodAllocationTrySuccess(nodeName string, pod *v1.Pod)  {
	refreshed, _ := lock.GetClient().CoreV1().Pods(pod.Namespace).Get(context.Background(), pod.Name, metav1.GetOptions{})
	
	annos := refreshed.Annotations[xpu.AssignedDevicesToAllocate]
	log.InfoIn("TrySuccess:", annos)
	
	if strings.Contains(annos, xpu.DeviceType) {
		return
	}

	log.Infoln("AllDevicesAllocateSuccess releasing lock")
	PodAllocationSuccess(nodeName, pod)
}

// PodAllocationSuccess patch annotation of a pod to indicate allocation success
func PodAllocationSuccess(nodeName string, pod *v1.Pod) {
	newAnnos := make(map[string]string)
	newAnnos[xpu.DeviceBindPhase] = types.DeviceBindSuccess
	err := PatchPodAnnotations(pod, newAnnos)
	if err != nil {
		log.Errorln("patchPodAnnotations failed:%v", err.Error())
	}
	err = lock.ReleaseNodeLock(nodeName, types.VXPULockName)
	if err != nil {
		log.Errorln("release lock failed:%v", err.Error())
	}
}

// PodAllocationFailed patch annotation of a pod to indicate allocation failed
func PodAllocationFailed(nodeName string, pod *v1.Pod) {
	newAnnos := make(map[string]string)
	newAnnos[xpu.DeviceBindPhase] = types.DeviceBindFailed
	err := PatchPodAnnotations(pod, newAnnos)
	if err != nil {
		log.Errorln("patchPodAnnotations failed:%v", err.Error())
	}
	err = lock.ReleaseNodeLock(nodeName, types.VXPULockName)
	if err != nil {
		log.Errorln("release lock failed:%v", err.Error())
	}
}

// PatchNodeAnnotations patch annotation of a node
func PatchNodeAnnotations(nodeName string, annotations map[string]string) error {
	type patchMetadata struct {
		Annotations map[string]string `json:"annotations,omitempty"`
	}
	type patchPod struct {
		Metadata patchMetadata `json:"metadata"`
	}
	p := patchPod{}
	p.Metadata.Annotations = annotations
	bytes, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = lock.GetClient().CoreV1().Nodes().Patch(context.Background(), nodeName, k8stypes.StrategicMergePatchType, bytes, metav1.PatchOptions{})
	if err != nil {
		log.Errorln("patch node %s failed: %v", nodeName, err)
	}
	return err
}

// PatchPodAnnotations patch annotation of a pod
func PatchPodAnnotations(pod *v1.Pod, annotations map[string]string) error {
	type patchMetadata struct {
		Annotations map[string]string `json:"annotations,omitempty"`
	}
	type patchPod struct {
		Metadata patchMetadata `json:"metadata"`
	}
	p := patchPod{}
	p.Metadata.Annotations = annotations
	bytes, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = lock.GetClient().CoreV1().Pods(pod.Namespace).Patch(context.Background(), pod.Name, k8stypes.StrategicMergePatchType, bytes, metav1.PatchOptions{})
	if err != nil {
		log.Errorln("patch pod %s failed: %v", pod.Name, err)
	}
	return err
}

// GetXpus description get xpu info on the node
func GetXpus() (types.VxpuDevices, error) {
	node, err := GetNode(config.NodeName)
	if err != nil {
		log.Errorln("get node error: %v, node Name: %s", err, config.NodeName)
		return nil, err
	}
	annos, ok := node.ObjectMeta.Annotations[xpu.NodeVGPURegister]
	if !ok {
		errMsg := fmt.Sprintf("node %s annotation %s is not exists",
			config.NodeName, xpu.NodeVGPURegister)
		log.Errorln(errMsg)
		return nil, errors.New(errMsg)
	}
	ips := getNodeIP(node)
	return GetXpuDevice(annos, ips), nil
}

func getNodeIP(node *v1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == v1.NodeInternalIP {
			return addr.Address
		}
	}
	log.Infoln("no expected node ip")
	return ""
}

// GetVgpus get all the xpu device info of the node
func GetVgpus() (types.VxpuDevices, map[string]uint32, error) {
	podList, err := ListPods(metav1.ListOptions{
		FieldSelector: selector.String(),
	})
	if err != nil {
		log.Errorln("get pods in current node error: %v", err)
		return nil, nil, err
	}
	res := types.VxpuDevices{}
	post := make(map[string]uint32)
	for _, pod := range podList.Items {
		if podStatus.Phase != v1.PodRunning || len(podStatus.ContainerStatuses) == 0 {
			log.Infoln("pod %s phase: %v, container status len: %v", pod.UID,
				podStatus.Phase, len(podStatus.ContainerStatuses))
			continue
		}
		pdevices := DecodePodDevices(pod.Annotations[xpu.AssignedIDs])
		pi := 0
		for _, cs := range pod.Spec.Containers {
			number, core, mem := getVgpuLimit(cs.Resources.Limits)
			// If the container has not configured resources.limits
			// it means that the container has no vgpu.
			if number == 0 {
				continue
			}
			// check the length of pdevices
			if pi >= len(pdevices) {
				log.Errorln("pod %v does not have enough xpu devices for %v", pod.UID, pi)
				break
			}
			// The vgpu number should be equal to the length of types.ContainerDevices.
			if len(devices[pi]) != int(number) {
				log.Warningf("xpu assigned info error, pod UID: %v, container name: %s", pod.UID, cs.Name)
				continue
			}
			for i := 0; i < int(number); i++ {
				dev := types.VxpuDevice{
					Id:                    fmt.Sprintf("%s-%d", pdevices[pi][i].UUID, pdevices[pi][i].Vid),
					GpuId:                 pdevices[pi][i].UUID,
					PodUID:                string(pod.UID),
					ContainerName:         cs.Name,
					VxpuMemoryLimit:       mem * 1024,
					VxpuCoreLimit:         core,
				}
				pi += 1
				key := fmt.Sprintf("%s/%s", string(pod.UID), cs.Name)
				post[key] = []uint32{}
			}
		}
	}
	return res, post, nil
}