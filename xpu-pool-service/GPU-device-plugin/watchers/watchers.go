/*
Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
*/

// Package watchers create FSWatcher and OsWatcher
package watchers

import (
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"
)

// NewFSWatcher new file watcher
func NewFSWatcher(files ...string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if err = watcher.Add(f); err != nil {
			watcher.Close()
			return nil, err
		}
	}

	return watcher, nil
}

// NewOSWatcher new os signal watcher
func NewOSWatcher(sigs ...os.Signal) chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, sigs...)
	return sigChan
}