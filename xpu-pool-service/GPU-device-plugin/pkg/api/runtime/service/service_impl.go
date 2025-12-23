/*
Copyright (c) Huawei Technologies Co., Ltd. 2024-2025. All rights reserved.
*/

// Package service implements service of getting pids
package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"

	"huawei.com/xpu-device-plugin/pkg/log"
	"huawei.com/xpu-device-plugin/pkg/plugin/config"
	"huawei.com/xpu-device-plugin/pkg/plugin/types"
	"huawei.com/xpu-device-plugin/pkg/plugin/util"
	"huawei.com/xpu-device-plugin/pkg/xpu"
)

const (
	cgroupBaseDir         = "/sys/fs/cgroup/memory"
	cgroupProcs           = "cgroup.procs"
	hostProcDir           = "/host/proc"
	hostProcStat          = "/var/lib/xpu/pids.sock"
	procStat              = "stat"
	npPid                 = "NpPid"
	npPidFieldCount       = 3
	containerPidDirPrefix = "/pod"
	containerPidSuffix    = ".pod"
	dockerPidPrefix       = "/docker"
	dockerPidSuffix       = ".docker"
	containerdIdPrefix    = "cri-containerd-"
	containerdIdSuffix    = ".scope"
	containerdPrefixInContainerd = "containerd://"
	containerdPrefixInDocker     = "docker://"
	vgpuConfigBaseDir    = "/etc/xpu"
	vgpuConfigFileName   = "pids.config"
	configFilePerm       = 0644
	pidsSockPerm         = 0666
	podsDirCleanInterval = 60
	minPeriod            = 1
	maxPeriod            = 86400
	defaultPeriod        = 60
	percentRange         = 100
	float64BitSize      = 64
)

// PidsServiceServerImpl implementation of pids service
type PidsServiceServerImpl struct {
	*UnimplementedPidsServiceServer
}

func readProcsFile(file string) ([]int, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Error("can't read %s: %s", file, err.Error())
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	pids := make([]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if pid, err := strconv.Atoi(line); err == nil {
			pids = append(pids, pid)
		}
	}

	log.Infof("read from %s, pids: %v", file, pids)
	return pids, nil
}

func readStatusFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Error("can't read %s", file, err.Error())
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, npPid) {
			eles := strings.Fields(line)
			if len(eles) != npPidFieldCount {
				return "", errors.New("NpPid field count error")
			}
			pids := fmt.Sprintf("%s %s", eles[1], eles[2])
			klog.Infof("read from %s, pids: %s", file, pids)
			return pids, nil
		}
	}
	return "", errors.New("NpPid not found")
}

func getHostPids() (string, error) {
	pidMaps := make([]string, 0)
	for _, hp := range hostPids {
		procStatusPath := filepath.Join(hostProcDir, strconv.Itoa(hp), procStat)
		pids, err := readStatusFile(procStatusPath)
		if err != nil {
			klog.Warning("read proc status error: %v, path: %s", err, procStatusPath)
		} else {
			pidMaps = append(pidMaps, pids)
		}
	}
	return strings.Join(pidMaps, ","), nil
}

func parseDockerCgroupPath(cgroupPath string) (string, string, error) {
	idx := strings.Index(cgroupPath, dockerPidPrefix)
	if idx == -1 {
		return "", "", errors.New("pod id prefix not found")
	}
	podIdTmp := cgroupPath[idx+len(dockerPidPrefix):]
	idx = strings.Index(podIdTmp, dockerPidSuffix)
	if idx == -1 {
		return "", "", errors.New("pod id suffix not found")
	}
	podId := podIdTmp[:idx]
	podId = strings.Replace(podId, ".", "-", -1)
	containerIdTmp := podIdTmp[idx+len(dockerPodIdSuffix):]
	containerId := containerIdTmp
	return podId, containerId, nil
}

func parseCgroupPath(cgroupPath string) (string, string, error) {
	idx := strings.Index(cgroupPath, containerPodIdPrefix)
	if idx == -1 {
		return parseDockerCgroupPath(cgroupPath)
	}
	podIdTmp := cgroupPath[idx+len(containerPodIdPrefix):]
	idx = strings.Index(podIdTmp, containerPodIdSuffix)
	if idx == -1 {
		return "", "", errors.New("pod id suffix not found")
	}
	podId := podIdTmp[:idx]
	podId = strings.Replace(podId, "_", "-", -1)
	idx = strings.Index(cgroupPath, containerIdPrefix)
	if idx == -1 {
		return podId, "", errors.New("container id prefix not found")
	}
	containerIdTmp := cgroupPath[idx+len(containerIdPrefix):]
	idx = strings.Index(containerIdTmp, containerIdSuffix)
	if idx == -1 {
		return podId, "", errors.New("container id suffix not found")
	}
	containerId := containerIdTmp[:idx]
	return podId, containerId, nil
}

func getContainerName(cgroupPath string) (string, string, error) {
	podId, containerId, err := parseCgroupPath(cgroupPath)
	if err != nil {
		log.Error("parse cgroup path error: %v", err)
		return "", "", err
	}
	selector := fields.SelectorFromSet(fields.Set{"spec.nodeName": config.NodeName, "status.phase": string(v1.PodRunning)})
	podList, err := util.ListPodsWithOptions(
		metav1.ListOptions{
			FieldSelector: selector.String(),
		})
	if err != nil {
		log.Error("get pods in current node error: %v", err)
		return podId, "", err
	}
	for _, pod := range podList.Items {
		if string(pod.UID) != podId {
			continue
		}
		if pod.Status.Phase != v1.PodRunning || len(pod.Status.ContainerStatuses) == 0 {
			errMsg := fmt.Sprintf("pod status error: %v, container status len: %d",
				pod.Status.Phase, len(pod.Status.ContainerStatuses))
			log.Error(errMsg)
			return podId, "", errors.New(errMsg)
		}
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.ContainerID[len(containerIdPrefixInContainerd):] != containerId &&
				cs.ContainerID[len(containerIdPrefixInDocker):] != containerId {
				continue
			}
			return podId, cs.Name, nil
		}
	}
	return podId, "", errors.New("container not found")
}

func readPidsConfig(pidsConfigPath string) ([]uint32, error) {
	f, err := os.OpenFile(pidsConfigPath, os.O_RDONLY, configFilePerm)
	if err != nil {
		log.Error("open pids config file error: %v", err)
		return nil, err
	}
	defer f.Close()

	var pids []uint32
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		tmp := strings.Fields(line)
		if len(tmp) == 0 {
			continue
		}
		pid, err := strconv.Atoi(tmp[0])
		if err != nil {
			log.Warning("pid is wrong: %v", line)
			continue
		}
		pids = append(pids, uint32(pid))
	}
	return pids, nil
}

func writePidsConfig(pidsConfigPath, pidMaps string) error {
	f, err := os.OpenFile(pidsConfigPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, configFilePerm)
	if err != nil {
		log.Error("open pids config file error: %v", err)
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	pids := strings.Split(pidMaps, ",")
	for _, pid := range pids {
		if err = w.WriteString(pid + "\n"); err != nil {
			log.Error("bufio Writer WriteString error: %v", err)
			return err
		}
	}
	return w.Flush()
}

func getPodDirNames() ([]string, error) {
	dirNames := make([]string, 0)
	err := filepath.Walk(vgpuConfigBaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error("filepath walk func error: %v", err)
			return err
		}
		if info == nil || info.IsDir() || path == vgpuConfigBaseDir {
			return nil
		}
		dirNames = append(dirNames, info.Name())
		return filepath.SkipDir
	})
	if err != nil {
		log.Error("filepath walk error: %v, dir name: %s", err, vgpuConfigBaseDir)
		return []string{}, err
	}
	return dirNames, nil
}

type void struct{}

var val void

func cleanDestroyedPodDir() error {
	podDirNames, err := getPodDirNames()
	if err != nil {
		return err
	}

	selector := fields.SelectorFromSet(fields.Set{"spec.nodeName": config.NodeName})
	podList, err := util.ListPodsWithOptions(
		metav1.ListOptions{
			FieldSelector: selector.String(),
		})
	if err != nil {
		log.Error("get pods in current node error: %v", err)
		return err
	}

	podIdSet := make(map[string]void)
	for _, pod := range podList.Items {
		podIdSet[string(pod.UID)] = val
	}
	for _, podDirName := range podDirNames {
		if _, ok := podIdSet[podDirName]; ok {
			continue
		}
		podAbsoluteDir := filepath.Clean(filepath.Join(vgpuConfigBaseDir, podDirName))
		err = os.RemoveAll(podAbsoluteDir)
		if err != nil {
			log.Info("remove pod dir error: %v, dir name: %s", err, podAbsoluteDir)
		}
	}
	return nil
}

// GetPids pids service external interface, get all pids map relationship in container
func (PidsServiceServerImpl) GetPids(ctx context.Context, req *GetPidsRequest) (*GetPidsResponse, error) {
	cgroupAbsolutePath := filepath.Clean(filepath.Join(cgroupBaseDir, req.CgroupPath, cgroupProcs))
	hostPids, err := readProcsFile(cgroupAbsolutePath)
	if err != nil {
		log.Error("read group procs error: %v, path: %s", err, cgroupAbsolutePath)
		return nil, err
	}
	pidMaps := getPidMaps(hostPids)
	podId, containerName, err := getContainerName(req.CgroupPath)
	if err != nil {
		log.Error("get container name error: %v, cgroup path: %s", err, req.CgroupPath)
		return nil, err
	}
	pidsConfigPath := filepath.Clean(filepath.Join(vgpuConfigBaseDir, podId, containerName, pidsConfigFileName))
	err = writePidsConfig(pidsConfigPath, pidMaps)
	if err != nil {
		log.Error("write pids config error: %v, podId: %s, containerName: %s", err, podId, containerName)
		return nil, err
	}
	return &GetPidsResponse{EncodedPids: pidMaps}, nil
}

func getPodSet(pSet map[string][]uint32) map[string][]uint32 {
	if pSet == nil {
		return nil
	}
	for k := range pSet {
		pidsConfigPath := filepath.Clean(filepath.Join(vgpuConfigBaseDir, k, pidsConfigFileName))
		pids, err := readPidsConfig(pidsConfigPath)
		if err != nil {
			log.Error("read pids config failed, err msg: %v", err)
			continue
		}
		pSet[k] = pids
	}
	return pSet
}

func setVgpuDevices(xpuDevices types.VgpuDevices,
	uidToProcessMap map[string]map[uint32]types.ProcessUsage,
	pSet map[string][]uint32) map[string]types.XpuDevice {
	for _, v := range xpuDevices {
		if _, ok := xpuDevices[v.Gpuid]; !ok {
			log.Warning("vgpu in pod %s container %s is not exist", v.PodUID, v.ContainerName)
			continue
		}
		// Get processUsage of xpu
		processUsage, ok := uidToProcessMap[v.Gpuid]
		if !ok {
			xpuDevices[v.Gpuid].VgpuDeviceList = append(xpuDevices[v.Gpuid].VgpuDeviceList, v)
			continue
		}
		// Get pidlist of the container
		// Then we get intersection of pidlist and xpu processUsage
		// take the pid usage of the container's vgpu on this xpu
		pidList, ok := pSet[v.PodUID]
		if !ok {
			continue
		}
		v.VgpuCoreUtilization += float64(processUsage.ProcessCoreUtilization)
		v.VgpuMemoryUsed += float64(processUsage.MemoryUsed)
		v.VgpuMemoryUsed = v.VgpuMemoryUsed / 1024 / 1024
		if xpuDevices[v.Gpuid].Memory.Total == 0 {
			xpuMemoryUtil := float64(v.VgpuMemoryUsed*percentage) / float64(xpuDevices[v.Gpuid].Memory.Total)
			err := nil
			v.VgpuMemoryUtilization, err = strconv.ParseFloat(fmt.Sprintf("%.2f", xpuMemoryUtil), float64BitSize)
			if err != nil {
				v.VxpuMemoryUtilization = xpuMemoryUtil
			}
		}
		xpuDevices[v.Gpuid].VgpuDeviceList = append(xpuDevices[v.Gpuid].VgpuDeviceList, v)
	}
	for _, device := range xpuDevices {
		for _, vxpuDevice := range device.VgpuDeviceList {
			memoryUtil := float64(vxpuDevice.VgpuMemoryUsed) / float64(device.Memory.Total)
			memoryUtil, err := strconv.ParseFloat(fmt.Sprintf("%.2f", memoryUtil), float64BitSize)
			if err != nil {
				log.Warning("vgpuMemoryUtil Parse failed %v", memoryUtil)
			}
			device.MemoryUtilization = memoryUtil
		}
	}
	return xpuDevices
}

// GetAllVxpuInfo get all vxpu info of the node
func (PidsServiceServerImpl) GetAllVgpuInfo(ctx context.Context, req *GetAllVgpuInfoRequest) (*GetAllVxpuInfoResponse, error) {
	getAllVgpuInfoResponse := &GetAllVgpuInfoResponse{}
	period, err := strconv.Atoi(req.Period)
	if err != nil || period < minPeriod || period > maxPeriod {
		log.Warning("period is unqualified, err msg: %v", req.Period, err)
		period = defaultPeriod
	}

	// xpuDevices: map[uuid]xpuDevice
	xpuDevices, err := util.GetXPU()
	if err != nil {
		log.Error("get xpu device info failed, err msg: %v", err)
		return nil, err
	}

	// vgpuDevices: types.VgpuDevices[]
	vgpuDevices, pSet, err := util.GetVgpus()
	if err != nil {
		log.Error("get vgpu device info failed, err msg: %v", err)
		return nil, err
	}

	// pSet: map["podId/containerName"]pids
	pSet = getPodSet(pSet)

	// uidToProcessMap: map[uuid]map[processId]processUsage
	uidToProcessMap := make(map[string]map[uint32]types.ProcessUsage)
	for _, v := range xpuDevices {
		deviceUsageInfo, processMap, err := xpu.GetXPUUsage(v.Index, int32(period))
		if err != nil {
			log.Error("get xpu usage failed: %v", err)
			return nil, err
		}
		// XpuUtilization = deviceUsageInfo.CoreUtil
		v.Power = deviceUsageInfo.Power
		v.Temperature = deviceUsageInfo.Temperature
		uidToProcessMap[v.UUID] = processMap
	}
	xpuDevices = setVgpuDevices(xpuDevices, vgpuDevices, uidToProcessMap, pSet)
	jsonVxpuInfos, err := json.Marshal(xpuDevices)
	if err != nil {
		log.Error("Generate JSON string error: %v", err)
		return nil, err
	}
	return &GetAllVxpuInfoResponse{VxpuInfos: string(jsonVxpuInfos)}, nil
}

// Start run pids service
func Start() {
	srv := grpc.NewServer()
	RegisterPidsServiceServer(srv, PidsServiceServerImpl{})
	err := syscall.Unlink(pidsSockPath)
	if err != nil && !os.IsNotExist(err) {
		log.Error("remove pids server sock error: %v", err)
		return
	}
	listener, err := net.Listen("unix", pidsSockPath)
	if err != nil {
		log.Error("net listen error: %v", err)
		return
	}
	err = os.Chmod(pidsSockPath, pidsSockPerm)
	if err != nil {
		log.Error("modify pids.sock file permissions error: %v", err)
		return
	}
	go func() {
		err := srv.Serve(listener)
		if err != nil {
			log.Error("grpc server error: %v", err)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second * podDirCleanInterval)
			err := cleanDestroyedPodDir()
			if err != nil {
				log.Error("clean destroyed pod dir error: %v, stop cleanup...", err)
				break
			}
		}
	}()
}
