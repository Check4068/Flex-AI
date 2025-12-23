/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */
#include <fcntl.h>
#include "log.h"
#include "memory_limiter.h"

bool MemoryLimiter::MemoryCheck(size_t requested)
{
    if (!config_.LimitMemory()) {
        return true;
    }

    size_t used;
    int ret = xpu_.MemoryUsed(used);
    if (ret) {
        log_err("get used memory failed, ret is {%d}", ret);
        return false;
    }

    size_t quota = config_.MemoryQuota();
    if (requested + used > quota) {
        log_err("out of memory, request {%lld}, used {%lld}, quota {%lld}",
            requested, used, quota);
        return false;
    }
    return true;
}

MemoryLimiter::Guard MemoryLimiter::GuardedMemoryCheck(size_t requested)
{
    FileLock lock(LockPath(), LOCK_EX);
    return {std::move(lock), MemoryCheck(requested)};
}

int MemoryLimiter::CreateFileLockBaseDir()
{
    int ret = mkdir(FILELOCK_BASE_DIR.c_str(), S_IRWXU | S_IRGRP | S_IXGRP);
    if (ret < 0 && errno != EEXIST) {
        log_err("mkdir {%s} failed, err is {%d}", FILELOCK_BASE_DIR, strerror(errno));
        return RET_FAIL;
    }
    log_info("mkdir {%s} succ", FILELOCK_BASE_DIR);
    return RET_SUCC;
}

int MemoryLimiter::Initialize()
{
    return CreateFileLockBaseDir();
}