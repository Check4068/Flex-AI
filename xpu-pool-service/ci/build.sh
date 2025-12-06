#!/bin/bash
# Copyright (C) Huawei Technologies Co., Ltd. 2024-2025. All rights reserved.
set -e

WORK_DIR=$(cd "$(dirname "$0")" && pwd)
DEST_DIR=${WORK_DIR}/../xpu_pool/xpu_docker_build/

function prepare() {
    mkdir -p ${DEST_DIR}/cuda_client/GPU_client/
    mkdir -p ${DEST_DIR}/acl_client/NPU_client/
    mkdir -p ${WORK_DIR}/../XPU_symbols/
}

function handle_spdlog() {
    mkdir -m 750 -p /usr/local/include
    cd third_party/spdlog
    cp -P --remove-destination -rf include/spdlog /usr/local/include
    chmod 750 -R /usr/local/include
    cd ${WORK_DIR}
}

function compile_client() {
    cd ${WORK_DIR} && rm -rf build && mkdir build && cd build
    cmake -DCMAKE_BUILD_TYPE=Release ../.. && make -j
}

function strip_gotest_codes() {
    if [ ! -d "$1" ]; then
        echo "Error: Directory '$1' does not exist."
        return
    fi
    cd "$1"
    find . -name *.test.go | xargs rm -rf
    sed -i '/gomonkey/d' go.mod
    go mod tidy
}

function compile_device_plugin() {
    # strip gomoney related codes to make SwInfoTree happy
    strip_gotest_codes "${WORK_DIR}/../GPU-device-plugin/"
    cd ${WORK_DIR}/../GPU-device-plugin/ && make -j
}

function compile_xpu_exporter() {
    # strip gomoney related codes to make SwInfoTree happy
    strip_gotest_codes "${WORK_DIR}/../xpu-exporter/"
    cd ${WORK_DIR}/../xpu-exporter/ && make clean && make -j
}

function strip_symbols() {
    cd ${WORK_DIR}/build/direct/cuda
    objcopy --only-keep-debug libcuda_direct.so libcuda_direct.sym
    objcopy --only-keep-debug gpu-monitor_direct.so gpu-monitor_direct.sym
    objcopy --strip-all libcuda_direct.so
    objcopy --strip-all gpu-monitor_direct.so

    cd ${WORK_DIR}/build/direct/acl
    objcopy --only-keep-debug libruntime_direct.so libruntime_direct.sym
    objcopy --only-keep-debug npu-monitor_direct.so npu-monitor_direct.sym
    objcopy --strip-all libruntime_direct.so
    objcopy --strip-all npu-monitor_direct.so
}

function copy_to_build_dir() {
    cd ${WORK_DIR}/build/direct/cuda
    cp -P --remove-destination -r direct/cuda/libcuda_direct.so ${DEST_DIR}/cuda_client/GPU_client/
    cp -P --remove-destination -r direct/cuda/gpu-monitor_direct.so ${DEST_DIR}/cuda_client/GPU_client/
    cp -P --remove-destination -r direct/cuda/libcuda_direct.sym ${DEST_DIR}/cuda_client/GPU_client/
    cp -P --remove-destination -r direct/cuda/gpu-monitor_direct.sym ${DEST_DIR}/cuda_client/GPU_client/

    cp -P --remove-destination -r direct/cuda/client_update/cuda_client-update.sh ${DEST_DIR}/cuda_client/GPU_client/

    cd ${WORK_DIR}/build/direct/acl
    cp -P --remove-destination -r direct/acl/libruntime_direct.so ${DEST_DIR}/acl_client/NPU_client/
    cp -P --remove-destination -r direct/acl/npu-monitor_direct.so ${DEST_DIR}/acl_client/NPU_client/
    cp -P --remove-destination -r direct/acl/libruntime_direct.sym ${DEST_DIR}/acl_client/NPU_client/
    cp -P --remove-destination -r direct/acl/npu-monitor_direct.sym ${DEST_DIR}/acl_client/NPU_client/

    cp -P --remove-destination -r direct/acl/client_update/acl_client-update.sh ${DEST_DIR}/acl_client/NPU_client/

    cd ${WORK_DIR}/../GPU-device-plugin
    cp -P --remove-destination -r gpu-device-plugin ${DEST_DIR}/gpu-device-plugin

    cd ${WORK_DIR}/../xpu-exporter
    cp -P --remove-destination -r xpu-exporter ${DEST_DIR}/exporter

    cd ${WORK_DIR}/../XPU_symbols
    tar -czvf XPU_symbols.tar.gz XPU_symbols/
    cp XPU_symbols.tar.gz ${DEST_DIR}/
}

function main() {
    prepare
    handle_spdlog
    compile_client
    compile_device_plugin
    compile_xpu_exporter
    strip_symbols
    copy_to_build_dir
}

main "$@"