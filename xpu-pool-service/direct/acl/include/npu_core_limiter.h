/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */
#ifndef NPU_CORE_LIMITER_H
#define NPU_CORE_LIMITER_H

#include <cstddef>
#include <mutex>
#include <condition_variable>
#include <atomic>
#include "common.h"
#include "npu_timeslice_scheduler.h"
#include "sem.h"
#include "shm.h"
#include "stream_cache.h"
#include "resource_config.h"
#include "npu_manager.h"

class NpuCoreLimiter {
public:
    class RequestGuard {
    public:
        friend class NpuCoreLimiter;
        RequestGuard(const RequestGuard& other) = delete;
        RequestGuard& operator=(const RequestGuard& other) = delete;

    private:
        explicit RequestGuard(NpuCoreLimiter &limiter, rtContext_t ctx, rtStream_t stream);
        NpuCoreLimiter &limiter_;
    };

    class ReleaseGuard {
    public:
        ReleaseGuard(const ReleaseGuard& other) = delete;
        ReleaseGuard& operator=(const ReleaseGuard& other) = delete;

    private:
        explicit ReleaseGuard(NpuCoreLimiter &limiter, size_t opCount);
        NpuCoreLimiter &limiter_;
        size_t opCount_;
    };

    RequestGuard ComputingPowerLimiter(rtStream_t stm);
    ReleaseGuard ReleaseOps(size_t ops);

    static NpuCoreLimiter Instance();

    NpuCoreLimiter(ResourceConfig &config, NpuManager &npu) : config_(config), npu_(npu)
    {}
    NpuCoreLimiter();

    void ComputingPowerWatcherThread();

    int Initialize();

TESTABLE_PROTECTED:
    int LoadNpuIdsConfig();
    int ComputingPowerWatcherInit();
    void JoinWatcher();

    ResourceConfig &config_;
    NpuManager &npu_;
    bool watcherEnd_ = false;

TESTABLE_PRIVATE:
    // UPDATE_PERIOD is an empirical value: 1/6 s
    constexpr static auto UPDATE_PERIOD = std::chrono::milliseconds(167);
    const std::string WUPLIDS_CONFIG_PATH = npu_CONFIG_BASE_DIR + "npu-ids.config";

    std::string strId_;
    int idx_ = -1;
    int32_t deviceCnt_ = XpuManager::INVALID_DEVICE_COUNT;
    Sem semaphore_;
    Sem semaphoreBack_;
    StreamCache streams_;
    std::thread watcher_;
    // Claim shared memory before scheduler, so that it will deconstruct scheduler before shared memory
    Shm shm_;
    NpuTimesliceScheduler sched_;
};

#endif