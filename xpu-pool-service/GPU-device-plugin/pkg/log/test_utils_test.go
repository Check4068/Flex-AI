/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

// Package log provides a rolling FileLogger.
package log

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func assertIsNil(obtained interface{}, t testing.TB) {
	isNil(obtained, t, 1)
}

func assertNotNil(obtained interface{}, t testing.TB) {
	notNil(obtained, t, 1)
}

func assertEquals(expect, act interface{}, t testing.TB) {
	equals(expect, act, t, 1)
}

func isNil(obj interface{}, t testing.TB, caller int) {
	if !_isNil(obj) {
		_, file, line, _ := runtime.Caller(caller + 1)
		t.Fatalf("%s:%d: expected nil, got: %#v", filepath.Base(file), line, obj)
		t.FailNow()
	}
}

func notNil(obj interface{}, t testing.TB, caller int) {
	if _isNil(obj) {
		_, file, line, _ := runtime.Caller(caller + 1)
		t.Fatalf("%s:%d: expected non-nil, got: %#v", filepath.Base(file), line, obj)
		t.FailNow()
	}
}

func _isNil(obj interface{}) bool {
	if obj == nil {
		return true
	}

	switch v := reflect.ValueOf(obj); v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func equals(expect, act interface{}, t testing.TB, caller int) {
	if !reflect.DeepEqual(expect, act) {
		_, file, line, _ := runtime.Caller(caller + 1)
		t.Fatalf("%s:%d: expected %#v, got %#v", filepath.Base(file), line, expect, act)
		t.FailNow()
	}
}