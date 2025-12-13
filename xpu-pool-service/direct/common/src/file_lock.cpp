/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */
#include <string>
#include <sys/file.h>
#include "log.h"
#include "common.h"
#include "file_lock.h"

FileLock::FileLock(const std::string_view path, int operation) : held_(false)
{
    fd_ = open(std::string(path).data(), O_CREAT | O_RDONLY, 0600); // the perm of lock file is 0600
    if (fd_ == -1) {
        log_err("open {%s} failed, errno is {%d}, %s", path, strerror(errno));
        return;
    }
    Acquire(operation);
}

bool FileLock::Acquire(int operation)
{
    /*
    * (1) The flock can block when anyone else holding the lock;
    * (2) The flock can be released either by calling the LOCK_UN parameter or
    *     by closing fd (the first parameter in flock), which means that flock
    *     is automatically released when the process is exited.
    */
    int ret = flock(fd_, operation);
    if (ret) {
        log_err("flock failed, fd {%d}, errno {%d}, %s", fd_, strerror(errno));
        return false;
    }
    held_ = true;
    return true;
}

bool FileLock::Release()
{
    int ret = flock(fd_, LOCK_UN);
    if (ret) {
        log_err("unlock failed, fd {%d}, errno {%d}, %s", fd_, strerror(errno));
        return false;
    }
    held_ = false;
    return true;
}

FileLock::~FileLock()
{
    if (fd_ < 0) {
        return;
    }
    if (held_) {
        Release();
    }
    if (close(fd_) == -1) {
        log_err("close file failed, fd {%d}, errno is {%d}, %s", fd_, strerror(errno));
    }
}