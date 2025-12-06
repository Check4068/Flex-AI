/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

// Package client implement grpc interface call to query vgpu information
package client

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"huawei.com/xpu-device-plugin/pkg/log"
	"huawei.com/xpu-exporter/common/service"
)

const (
	pidsSockPath = "/var/lib/xpu/pids.sock"
	dialTimeout  = 5
	megabyte     = 1024 * 1024
)

// GetAllVxpuInfo Obtain vgpu information through grpc interface
func GetAllVxpuInfo() (string, error) {
	conn, err := grpc.Dial(pidsSockPath,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(dialTimeout*time.Second),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		log.Errorf("grpc dial error: %v", err)
		return "", err
	}
	defer conn.Close()

	client := service.NewPidsServiceClient(conn)
	getAllVxpuInfoResponse, err := client.GetAllVxpuInfo(context.Background(), &service.GetAllVxpuInfoRequest{Period: "60"})
	if err != nil {
		log.Errorf("client GetAllVxpuInfo error: %v", err)
		return "", err
	}
	return getAllVxpuInfoResponse.VxpuInfos, nil
}