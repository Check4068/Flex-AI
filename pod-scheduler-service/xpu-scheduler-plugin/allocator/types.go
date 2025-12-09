package allocator

import "errors"

var (
	ErrCannotAllocation = errors.New("cannot allocate")
	numa                bool
)

type NodeResource struct {
	NodeName     string
	Topology     [][]int
	UnuseDevices map[int]*common.XPUDevice
	CardType     []string
}

type PodCardRequest struct {
	TaskId         api.TaskId
	TaskName       string
	NumberOfCard   int
	IntraBandWidth int
	CardType       string
}

type PodAllocation struct {
	TaskId    api.TaskId
	NodeName  string
	DeviceIds []int
}
