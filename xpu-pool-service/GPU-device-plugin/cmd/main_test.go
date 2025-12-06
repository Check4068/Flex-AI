// Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.

// Package main
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/agiledragon/gommonkey/v2"
	"github.com/fsnotify/fsnotify"
	"k8s.io/kubelet/device-plugin/vibetal"
	"huawei.com/xpu-device-plugin/pkg/api/runtime/service"
	"huawei.com/xpu-device-plugin/pkg/govmli"
	"huawei.com/xpu-device-plugin/pkg/plugin"
	"huawei.com/xpu-device-plugin/pkg/xpu"
	"huawei.com/xpu-device-plugin/watchers"
)

func ApplyPatches(patches []func() *gommonkey.Patches) []*gommonkey.Patches {
	var appliedPatches []*gommonkey.Patches
	for _, f := range patches {
		ap := f()
		appliedPatches = append(appliedPatches, ap)
	}
	return appliedPatches
}

func ResetPatches(appliedPatches []*gommonkey.Patches) {
	for _, p := range appliedPatches {
		p.Reset()
	}
}

var testCasesForEventWatcher = []struct {
	desc     string
	signal   string // if "panic" - expected function will enter this state
	expected bool
}{
	{
		desc:     "test for checking fsnotify event",
		signal:   "fsnotify",
		expected: true,
	},
	{
		desc:     "test for checking sighup event",
		signal:   "sighup",
		expected: true,
	},
	{
		desc:     "test for checking sigusr event",
		signal:   "sigusr",
		expected: false,
	},
}

func TestEventWatcher(t *testing.T) {
	for i := range testCasesForEventWatcher {
		tc := testCasesForEventWatcher[i]
		t.Run(tc.desc, func(t *testing.T) {
			w, _ := fsnotify.NewWatcher()
			s := make(chan os.Signal)
			pi := &plugin.DevicePlugin{}
			var res bool
			eventDoneCh := make(chan bool)
			go func() {
				res = events(w, s, pi)
				eventDoneCh <- true
			}()
			if tc.signal == "fsnotify" {
				w.Events <- fsnotify.Event{Name: vibetal.KubeletSocket, Op: fsnotify.Create}
			} else if tc.signal == "sighup" {
				s <- syscall.SIGHUP
			} else {
				s <- syscall.SIGUSR1
			}
			select {
			case <-eventDoneCh:
			case <-time.After(2 * time.Second): // just time for test case
				t.Error("event loop should exit, but timeout")
				return
			}
			if res != tc.expected {
				t.Error("event loop should with true status")
			}
		})
	}
}

var testCasesForStart = []struct {
	desc           string
	signal         string // if "panic" - expected function will enter this state
	patches        []func() *gommonkey.Patches
	expectedError  bool
}{
	{
		desc: "test for error handling from govmli.Init()",
		patches: []func() *gommonkey.Patches{
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(govmli.Init, func() govmli.NvmliRetType {
					return govmli.ErrorUnknown
				})
			},
		},
		expectedError: true,
	},
	{
		desc: "test for error handling from NewFSWatcher",
		patches: []func() *gommonkey.Patches{
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(govmli.Init, func() govmli.NvmliRetType {
					return govmli.Success
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(watchers.NewFSWatcher, func(files ...string) (*fsnotify.Watcher, error) {
					return nil, fmt.Errorf("test error from fs watcher")
				})
			},
		},
		expectedError: true,
	},
	{
		desc: "check error handling if devices = 0",
		patches: []func() *gommonkey.Patches{
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(govmli.Init, func() govmli.NvmliRetType {
					return govmli.Success
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(log.Info, func(format string, args ...interface{}) {
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(watchers.NewFSWatcher, func(files ...string) (*fsnotify.Watcher, error) {
					return nil, nil
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(plugin.NewDeviceCache, func() *plugin.DeviceCache {
					return &plugin.DeviceCache{}
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(plugin.NewDeviceRegister, func(c *plugin.DeviceCache) *plugin.DeviceRegister {
					return &plugin.DeviceRegister{}
				})
			},
			func() *gommonkey.Patches {
				var dc *plugin.DeviceCache
				return gommonkey.ApplyMethod(dc, "Start", func(*plugin.DeviceCache) {
				})
			},
			func() *gommonkey.Patches {
				var dc *plugin.DeviceCache
				return gommonkey.ApplyMethod(dc, "Stop", func(*plugin.DeviceCache) {
				})
			},
			func() *gommonkey.Patches {
				var dr *plugin.DeviceRegister
				return gommonkey.ApplyMethod(dr, "Start", func(*plugin.DeviceRegister) {
				})
			},
			func() *gommonkey.Patches {
				var w *fsnotify.Watcher
				return gommonkey.ApplyMethod(w, "Close", func(*fsnotify.Watcher) error {
					return nil
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(plugin.NewDevicePlugin, func(resourceName string, deviceCache *plugin.DeviceCache, socket string) *plugin.DevicePlugin {
					return &plugin.DevicePlugin{}
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(service.Start, func() {
				})
			},
			func() *gommonkey.Patches {
				var dr *plugin.DevicePlugin
				return gommonkey.ApplyMethod(dr, "Devices", func(*plugin.DevicePlugin) []*xpu.Device {
					return []*xpu.Device{}
				})
			},
		},
		expectedError: true,
	},
	{
		desc: "check error handling if start plugin error",
		patches: []func() *gommonkey.Patches{
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(govmli.Init, func() govmli.NvmliRetType {
					return govmli.Success
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(log.Info, func(format string, args ...interface{}) {
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(watchers.NewFSWatcher, func(files ...string) (*fsnotify.Watcher, error) {
					return nil, nil
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(plugin.NewDeviceCache, func() *plugin.DeviceCache {
					return &plugin.DeviceCache{}
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(plugin.NewDeviceRegister, func(c *plugin.DeviceCache) *plugin.DeviceRegister {
					return &plugin.DeviceRegister{}
				})
			},
			func() *gommonkey.Patches {
				var dc *plugin.DeviceCache
				return gommonkey.ApplyMethod(dc, "Start", func(*plugin.DeviceCache) {
				})
			},
			func() *gommonkey.Patches {
				var dc *plugin.DeviceCache
				return gommonkey.ApplyMethod(dc, "Stop", func(*plugin.DeviceCache) {
				})
			},
			func() *gommonkey.Patches {
				var dr *plugin.DeviceRegister
				return gommonkey.ApplyMethod(dr, "Start", func(*plugin.DeviceRegister) {
				})
			},
			func() *gommonkey.Patches {
				var w *fsnotify.Watcher
				return gommonkey.ApplyMethod(w, "Close", func(*fsnotify.Watcher) error {
					return nil
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(plugin.NewDevicePlugin, func(resourceName string, deviceCache *plugin.DeviceCache, socket string) *plugin.DevicePlugin {
					return &plugin.DevicePlugin{}
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(service.Start, func() {
				})
			},
			func() *gommonkey.Patches {
				var dr *plugin.DevicePlugin
				return gommonkey.ApplyMethod(dr, "Devices", func(*plugin.DevicePlugin) []*xpu.Device {
					return []*xpu.Device{}
				})
			},
			func() *gommonkey.Patches {
				var dr *plugin.DevicePlugin
				return gommonkey.ApplyMethod(dr, "Start", func(*plugin.DevicePlugin) error {
					return fmt.Errorf("test error from device plugin start")
				})
			},
		},
		expectedError: true,
	},
	{
		desc: "check error handling if start plugin error",
		patches: []func() *gommonkey.Patches{
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(govmli.Init, func() govmli.NvmliRetType {
					return govmli.Success
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(log.Info, func(format string, args ...interface{}) {
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(watchers.NewFSWatcher, func(files ...string) (*fsnotify.Watcher, error) {
					return nil, nil
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(plugin.NewDeviceCache, func() *plugin.DeviceCache {
					return &plugin.DeviceCache{}
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(plugin.NewDeviceRegister, func(c *plugin.DeviceCache) *plugin.DeviceRegister {
					return &plugin.DeviceRegister{}
				})
			},
			func() *gommonkey.Patches {
				var dc *plugin.DeviceCache
				return gommonkey.ApplyMethod(dc, "Start", func(*plugin.DeviceCache) {
				})
			},
			func() *gommonkey.Patches {
				var dc *plugin.DeviceCache
				return gommonkey.ApplyMethod(dc, "Stop", func(*plugin.DeviceCache) {
				})
			},
			func() *gommonkey.Patches {
				var dr *plugin.DeviceRegister
				return gommonkey.ApplyMethod(dr, "Start", func(*plugin.DeviceRegister) {
				})
			},
			func() *gommonkey.Patches {
				var w *fsnotify.Watcher
				return gommonkey.ApplyMethod(w, "Close", func(*fsnotify.Watcher) error {
					return nil
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(plugin.NewDevicePlugin, func(resourceName string, deviceCache *plugin.DeviceCache, socket string) *plugin.DevicePlugin {
					return &plugin.DevicePlugin{}
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(service.Start, func() {
				})
			},
			func() *gommonkey.Patches {
				var dr *plugin.DevicePlugin
				return gommonkey.ApplyMethod(dr, "Devices", func(*plugin.DevicePlugin) []*xpu.Device {
					return []*xpu.Device{}
				})
			},
			func() *gommonkey.Patches {
				var dr *plugin.DevicePlugin
				return gommonkey.ApplyMethod(dr, "Start", func(*plugin.DevicePlugin) error {
					return nil
				})
			},
			func() *gommonkey.Patches {
				return gommonkey.ApplyFunc(events, func(watcher *fsnotify.Watcher, sigs chan os.Signal, pluginInst *plugin.DevicePlugin) bool {
					return false
				})
			},
		},
		expectedError: false,
	},
}

func TestStart(t *testing.T) {
	for i := range testCasesForStart {
		tc := testCasesForStart[i]
		t.Run(tc.desc, func(t *testing.T) {
			appliedPatches := ApplyPatches(tc.patches)
			defer ResetPatches(appliedPatches)
			err := Start()
			if err != nil && !tc.expectedError {
				t.Error("error in test case: ", tc.desc)
			} else if err == nil && tc.expectedError {
				t.Error("error in test case: ", tc.desc)
			}
		})
	}
}