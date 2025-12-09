package plugin

import (
	"fmt"
	"strconv"
	"strings"
)

func EncodeNodeDevices(xpuDevices []*common.XPUDevice) string {
	var encodeNodeDevices strings.Builder
	for _, val := range xpuDevices {
		encodeNodeDevices.Write([]byte(strconv.Itoa(val.Index)))
		encodeNodeDevices.Write([]byte(","))
		encodeNodeDevices.Write([]byte(val.Id))
		encodeNodeDevices.Write([]byte(","))
		encodeNodeDevices.Write([]byte(strconv.Itoa(val.Count)))
		encodeNodeDevices.Write([]byte(","))
		encodeNodeDevices.Write([]byte(strconv.Itoa(int(val.Memory))))
		encodeNodeDevices.Write([]byte(","))
		encodeNodeDevices.Write([]byte(val.Type))
		encodeNodeDevices.Write([]byte(","))
		encodeNodeDevices.Write([]byte(strconv.FormatBool(val.Health)))
		encodeNodeDevices.Write([]byte(":"))
	}
	return encodeNodeDevices.String()
}

func DecodeNodeDevices(str string, nodeId string) map[int]*common.XpuDevice {
	xpuDevices := make(map[int]*common.XpuDevice)
	if !strings.Contains(str, ":") {
		return xpuDevices
	}
	tmp := strings.Split(str, ":")
	for _, val := range tmp {
		if strings.Contains(val, ",") {
			items := strings.Split(val, ",")
			if len(items) != util.XPUDeviceLen {
				return map[int]*common.XpuDevice{}
			}
			index, err := strconv.Atoi(items[0])
			count, err := strconv.Atoi(items[2])
			memory, err := strconv.Atoi(items[3])
			health, err := strconv.ParseBool(items[5])
			numa, err := strconv.Atoi(items[6])
			if err != nil {
				return map[int]*common.XpuDevice{}
			}
			i := &common.XPUDevice{
				Index:      index,
				Id:         items[1],
				NodeId:     nodeId,
				Type:       items[4],
				Count:      count,
				Health:     health,
				Cores:      util.Base100,
				Memory:     uint64(memory),
				UsedCores:  0,
				UsedMemory: 0,
				UsedVids:   0,
				InUse:      false,
				Numa:       numa,
			}
			xpuDevices[index] = i
		}
	}
	return xpuDevices
}

func EncodeContainerDevices(cd ContainerDevices) string {
	var encodeContainerDevices strings.Builder
	for _, val := range cd {
		encodeContainerDevices.Write([]byte(strconv.Itoa(int(val.Index))))
		encodeContainerDevices.Write([]byte(","))
		encodeContainerDevices.Write([]byte(val.Id))
		encodeContainerDevices.Write([]byte(","))
		valType := val.Type
		if strings.Contains(valType, util.NvidiaGPUDevice) {
			valType = util.NvidiaGPUDevice
		}
		if strings.Contains(valType, util.AscendNPUDevice) {
			valType = util.AscendNPUDevice
		}
		encodeContainerDevices.Write([]byte(valType))
		encodeContainerDevices.Write([]byte(","))
		encodeContainerDevices.Write([]byte(strconv.Itoa(int(val.UsedMemory))))
		encodeContainerDevices.Write([]byte(","))
		encodeContainerDevices.Write([]byte(strconv.Itoa(int(val.UsedCores))))
		encodeContainerDevices.Write([]byte(","))
		encodeContainerDevices.Write([]byte(strconv.FormatUint(uint64(val.Vod), util.Base10)))
		encodeContainerDevices.Write([]byte(":"))
	}
	return (EncodeContainerDevices)
}
