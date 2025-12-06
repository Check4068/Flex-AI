/*
Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
*/

package watchers

import (
	"os"
	"syscall"
	"testing"
)

func TestInvalidNewFSWatcher(t *testing.T) {
	errTestCases := []struct {
		fileNames []string
	}{
		{fileNames: []string{"/non-existent-path/test-file-name"}},
	}

	for i := range errTestCases {
		tc := errTestCases[i]
		_, err := NewFSWatcher(tc.fileNames...)
		if err == nil {
			t.Errorf("no error on invalid filenames test case: %d", i)
		}
	}
}

func TestValidNewFSWatcher(t *testing.T) {
	testFileName := "/tmp/testfile"
	f, err := os.Create(testFileName)
	if err != nil {
		t.Error("cannot create file for test")
	}
	defer f.Close()

	defer func() {
		err = os.Chmod(testFileName, 0600) // file permissions
		if err != nil {
			t.Error("cannot change permissions for valid file")
		}
	}()

	_, err = NewFSWatcher(testFileName)
	if err != nil {
		t.Error("cannot create watcher for valid file")
	}

	err = os.Remove(testFileName)
	if err != nil {
		t.Error("cannot delete testfile after test")
	}
}

func TestNewOSWatcher(t *testing.T) {
	errTestCases := []struct {
		sigs []os.Signal
	}{
		{sigs: []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}},
		{sigs: []os.Signal{syscall.SIGHUP}},
	}

	for i := range errTestCases {
		tc := errTestCases[i]
		_ = NewOSWatcher(tc.sigs...)
		// 注：原代码此处可能漏写了断言逻辑（截图中无报错处理）
	}
}