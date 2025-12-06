/*
 *Copyright(c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

// Package main implements xpu client tool
// 测试文件：用于测试 client.go 中的 updatePidsConfig 函数
package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"google.golang.org/grpc"

	"huawei.com/vxpu-device-plugin/pkg/api/runtime/service"
)

// mockClient 模拟 gRPC 客户端，用于测试场景中替换真实的 PidsServiceClient
type mockClient struct {
	retVal      error                          // 模拟返回的错误值
	resp        *service.GetPidsResponse       // GetPids 方法的模拟响应
	vxpuInfoResp *service.GetVxpuInfoResponse  // GetAllVxpuInfo 方法的模拟响应
}

// ApplyPatches 批量应用 gomonkey patches，用于在测试中替换函数行为
func ApplyPatches(patches []func() *gomonkey.Patches) []*gomonkey.Patches {
	var appliedPatches []*gomonkey.Patches
	// 遍历所有 patch 函数，执行并收集返回的 Patches 对象
	for _, f := range patches {
		ap := f()
		appliedPatches = append(appliedPatches, ap)
	}
	return appliedPatches
}

// ResetPatches 清理所有已应用的 patches，恢复原始函数行为
func ResetPatches(appliedPatches []*gomonkey.Patches) {
	// 遍历所有 patches，逐个恢复
	for _, p := range appliedPatches {
		p.Reset()
	}
}

// GetPids 实现 PidServiceClient 接口的 GetPids 方法，返回预设的响应和错误
func (mc *mockClient) GetPids(ctx context.Context, req *service.GetPidsRequest, opts ...grpc.CallOption) (*service.GetPidsResponse, error) {
	return mc.resp, mc.retVal
}

// GetAllVxpuInfo 实现 VxpuInfoServiceClient 接口的 GetAllVxpuInfo 方法，返回预设的响应和错误
func (mc *mockClient) GetAllVxpuInfo(ctx context.Context, req *service.GetVxpuInfoRequest, opts ...grpc.CallOption) (*service.GetAllVxpuInfoResponse, error) {
	return mc.vxpuInfoResp, mc.retVal
}

// registerTestCaseForUpdatePidsConfig 定义 updatePidsConfig 函数的表驱动测试用例
// 使用 gomonkey 进行函数替换，模拟不同的测试场景
var registerTestCaseForUpdatePidsConfig = []struct {
	desc     string                   // 测试用例描述
	patches  []func() *gomonkey.Patches // patch 函数列表，用于替换被测试函数依赖的外部函数
	expected int                      // 期望结果：0-函数应返回无错误，1-函数应返回错误
}{
	{
		// 测试场景1：模拟 grpc.Dial 连接失败，验证错误处理逻辑
		desc: "checking error handling from grpc.Dial()",
		patches: []func() *gomonkey.Patches{
			func() *gomonkey.Patches {
				// 替换 grpc.Dial 函数，模拟连接失败返回错误
				return gomonkey.ApplyFunc(grpc.Dial, func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
					err := errors.New("test error")
					return nil, err
				})
			},
		},
		expected: 1,
	},
	{
		// 测试场景2：模拟 GetPids 调用失败，验证客户端调用错误处理
		desc: "checking error handling from GetPids()",
		patches: []func() *gomonkey.Patches{
			func() *gomonkey.Patches {
				// 替换 grpc.Dial 函数，模拟连接成功
				return gomonkey.ApplyFunc(grpc.Dial, func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
					return &grpc.ClientConn{}, nil
				})
			},
			func() *gomonkey.Patches {
				// 替换 NewPidsServiceClient 函数，返回会返回错误的 mock 客户端
				return gomonkey.ApplyFunc(service.NewPidsServiceClient, func(c grpc.ClientConnInterface) service.PidServiceClient {
					return &mockClient{retVal: fmt.Errorf("test error from GetPids")}
				})
			},
		},
		expected: 1,
	},
	{
		// 测试场景3：正常流程，所有调用都成功，验证 happy path
		desc: "happy path",
		patches: []func() *gomonkey.Patches{
			func() *gomonkey.Patches {
				// 替换 grpc.Dial 函数，模拟连接成功返回有效连接
				return gomonkey.ApplyFunc(grpc.Dial, func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
					return &grpc.ClientConn{}, nil
				})
			},
			func() *gomonkey.Patches {
				// 替换 NewPidsServiceClient 函数，返回正常的 mock 客户端（无错误）
				return gomonkey.ApplyFunc(service.NewPidsServiceClient, func(c grpc.ClientConnInterface) service.PidServiceClient {
					return &mockClient{retVal: nil}
				})
			},
			func() *gomonkey.Patches {
				var c *grpc.ClientConn
				return gomonkey.ApplyMethod( c,	 "Close", func(_ *grpc.ClientConn) error {
					return nil
				})
			},
		},
		expected: 0,
	},
	{
		// 测试场景4：模拟 GetPids 返回错误，验证错误处理（与场景2类似但包含 Close 方法 mock）
		desc: "checking error if GetPids return error",
		patches: []func() *gomonkey.Patches{
			func() *gomonkey.Patches {
				// 替换 grpc.Dial 函数，模拟连接成功返回有效连接
				return gomonkey.ApplyFunc(grpc.Dial, func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
					return &grpc.ClientConn{}, nil
				})
			},
			func() *gomonkey.Patches {
				// 替换 NewPidsServiceClient 函数，返回会返回错误的 mock 客户端
				return gomonkey.ApplyFunc(service.NewPidsServiceClient, func(c grpc.ClientConnInterface) service.PidServiceClient {
					return &mockClient{retVal: fmt.Errorf("test error from GetPids")}
				})
			},
			func() *gomonkey.Patches {
				var c *grpc.ClientConn
				return gomonkey.ApplyMethod( c,	 "Close", func(_ *grpc.ClientConn) error {
					return nil
				})
			},
		},
		expected: 1,
	},
}

// TestUpdatePidsConfig 测试 updatePidsConfig 函数的所有场景
// 使用表驱动测试，遍历所有测试用例，验证函数的错误处理逻辑
func TestUpdatePidsConfig(t *testing.T) {
	// 遍历所有定义的测试用例
	for i := range registerTestCaseForUpdatePidsConfig {
		tc := registerTestCaseForUpdatePidsConfig[i]
		// 使用 t.Run 为每个测试用例创建子测试，便于隔离和识别
		t.Run(tc.desc, func(t *testing.T) {
			// 应用所有 patches，替换被测试函数依赖的外部函数
			appliedPatches := ApplyPatches(tc.patches)
			// 测试结束后清理 patches，恢复原始函数行为
			defer ResetPatches(appliedPatches)
			
			// 设置测试用的 cgroup 路径
			path := "test"
			// 调用被测试函数
			err := updatePidsConfig(path)
			
			// 根据是否有错误，将实际结果转换为状态码（0-无错误，1-有错误）
			var status int
			if err != nil {
				status = 1
			}
			
			// 验证实际状态码是否与期望值一致
			if status != tc.expected {
				t.Errorf("test case '%s' failed: expected status %d, got %d, error: %v", tc.desc, tc.expected, status, err)
			}
		})
	}
}

// TestMainError 测试 main 函数在 updatePidsConfig 返回错误时的退出逻辑
func TestMainError(t *testing.T) {
	exited := false
	// 替换 updatePidsConfig 函数，模拟返回错误
	patchUpdatePidConfig := gomonkey.ApplyFunc(updatePidsConfig, func(target string) error {
		return fmt.Errorf("test error")
	})
	// 替换 os.Exit 函数，捕获退出调用但不实际退出
	patchOsExit := gomonkey.ApplyFunc(os.exit, func(c int) {
		exit = true
		return
	})
	defer patchUpdatePidConfig.Reset()
	defer patchOsExit.Reset()
	
	// 调用 main 函数
	main()
	if !exited	{
		t.Error("main should exit")
	}
}
