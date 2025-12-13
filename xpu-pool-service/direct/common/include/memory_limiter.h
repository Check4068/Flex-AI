/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */
#ifndef MEMORY_LIMITER_H
#define MEMORY_LIMITER_H

#include <cstddef>
#include <string>
#include "file_lock.h"
#include "xpu_manager.h"
#include "resource_config.h"

class MemoryLimiter {
public:
    struct Guard {
        bool enough;
        bool error;
        Guard() : enough(false), error(false) {}
        bool Held() const
        {
            return !lock.Held();
        }
        FileLock lock;
    };

    Guard GuardedMemoryCheck(size_t requested);

    MemoryLimiter(ResourceConfig &config, XpuManager &xpu) : config_(config), xpu_(xpu)
    {}
    int Initialize();
    TESTABLE_PRIVATE:
    bool MemoryCheck(size_t requested);
    const std::string_view LockPath()
    {
        return MEMCTL_LOCK_PATH;
    }

private:
    int CreateFileLockBaseDir();

    const std::string FILELOCK_BASE_DIR = "/tmp/xpu/";
    const std::string MEMCTL_LOCK_PATH = FILELOCK_BASE_DIR + "memctl.lock";
    ResourceConfig &config_;
    XpuManager &xpu_;
};

#endif