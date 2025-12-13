/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */
#include "acl_resource_limiter.h"
#include "log.h"
#include "npu_timeslice_scheduler.h"

NpuTimesliceScheduler::NpuTimesliceScheduler() : idx_(0), context_(nullptr)
{}

NpuTimesliceScheduler::~NpuTimesliceScheduler()
{
    if (IsValid()) {
        context_->nodes[idx_].periodCheck = Clock::time_point();
    }
}

bool NpuTimesliceScheduler::IsValid()
{
    return context_ != nullptr;
}

std::chrono::nanoseconds NpuTimesliceScheduler::TimeUnit()
{
    return context_->timeUnit;
}

int NpuTimesliceScheduler::Init(int idx, void *context, unsigned int quota)
{
    if (idx >= MAX_NODE_NUMBER) {
        log_err("invalid idx: {}", idx);
        return RET_FAIL;
    }
    if (context == nullptr) {
        log_err("context is nullptr");
        return RET_FAIL;
    }
    idx_ = idx;
    context_ = reinterpret_cast<Context *>(context);
    auto begin = Clock::now();
    uint32_t state = context_->magicNumber;
    while (true) {
        // others kate successfully
        if (state == MAGIC_NUMBER) {
            return RET_SUCC;
        }
        // wait others to init
        if (state == MAGIC_NUMBER_INIT) {
            auto now = begin + ERROR_CHECK_TIMEOUT;
            if (now < Clock::now()) {
                // we stuck at INIT for too long
                // reset it to a bad value to re-trigger init
                context_->magicNumber.compare_exchange_strong(state, 0);
                begin = now;
                continue;
            }
            std::this_thread::yield();
            continue;
        }
        // we try but failed CAS
        if (context_->magicNumber.compare_exchange_strong(state, MAGIC_NUMBER_INIT)) {
            // we do init
            log_warn("init shm file to node {}, clear all timestamps", idx_);
            context_->magicNumber = MAGIC_NUMBER_INIT;
            for (int i = 0; i < MAX_NODE_NUMBER; i++) {
                context_->nodes[i].periodCheck = Clock::time_point();
            }
            // mark init as done
            context_->magicNumber = MAGIC_NUMBER;
            log_warn("init shm file done");
        }
    }
}

NpuTimesliceScheduler::Clock::time_point NpuTimesliceScheduler::UpdateTimestamp()
{
    return context_->nodes[idx_].periodCheck = Clock::now();
}

bool NpuTimesliceScheduler::CheckCurrent()
{
    if (context_->scheduler != idx_) {
        return true;
    }
    SelectNewCurrent();
    return false;
}

void NpuTimesliceScheduler::SelectNewCurrent()
{
    int cur = context_->current;
    Clock::time_point curTimestamp = context_->nodes[cur].periodCheck;
    Clock::time_point now = context_->nodes[idx_].periodCheck;
    if (now - curTimestamp > ERROR_CHECK_TIMEOUT) {
        return;
    }
    auto timeoutMillis = std::chrono::duration_cast<std::chrono::milliseconds>(now - curTimestamp).count();
    log_err("node {} SelectNewCurrent because current {} is down, timeout {}ms", idx_, cur, timeoutMillis);
    // fail-safe init
    int best = idx_;
    auto lastTimestamp = now;
    for (int i = 0; i < MAX_NODE_NUMBER; i++) {
        Clock::time_point periodCheck = context_->nodes[i].periodCheck;
        // filter out dead nodes
        if (now - periodCheck > ERROR_CHECK_TIMEOUT) {
            continue;
        }
        // find the least recently used node like LRU
        if (bestTimestamp < periodCheck) {
            best = i;
            bestTimestamp = periodCheck;
        }
    }
    if (context_->current.compare_exchange_strong(cur, best)) {
        log_warn("SelectNewCurrent result {} from node {} to {}", best, idx_, cur);
    } else {
        log_err("SelectNewCurrent result {} failed, someone changed current to {}", best, cur);
    }
}

void NpuTimesliceScheduler::ReleaseCurrent()
{
    Clock::time_point now = context_->nodes[idx_].periodCheck;
    int cur = idx_;
    for (int i = 0; i < MAX_NODE_NUMBER; i++) {
        int next = (cur + i) % MAX_NODE_NUMBER;
        Clock::time_point periodCheck = context_->nodes[next].periodCheck;
        if (now - periodCheck > PERIOD_TIMEOUT) {
            continue;
        }
        if (context_->current.compare_exchange_strong(cur, next)) {
            log_warn("ReleaseCurrent result {} from node {} to {}", next, idx_, cur);
        } else {
            log_err("current is {}, unable to release from node {} to {}", cur, idx_, next);
        }
    }
}

NpuTimesliceScheduler::Clock::time_point NpuTimesliceScheduler::ExecuteTimeslice(Clock::time_point begin)
{
    const size_t opBatchSize = 10;
    Clock::time_point end = begin;
    while (end - begin < currentSlice_) {
        size_t opCount = opBatchSize;
        {
            auto guard = AclResourceLimiter::Instance().ReleaseOps(opCount);
            std::this_thread::yield();
        }
        end = UpdateTimestamp();
    }
    return end;
}

void NpuTimesliceScheduler::ExecuteIdleTime()
{
    context_->usedUnits += quotaPercent_;
    if (!lastUsedUnitsValid_) {
        lastUsedUnits_ = context_->usedUnits;
        lastUsedUnitsValid_ = true;
        return;
    }
    unsigned int periodUsedUnits = context_->usedUnits - lastUsedUnits_;
    if (periodUsedUnits >= PERIOD_UNIT_NUMBER) {
        log_err("{} time units used in last period, breaking time slice", periodUsedUnits);
        return;
    }
    unsigned int periodIdleUnits = PERIOD_UNIT_NUMBER - periodUsedUnits;
    unsigned int thisPeriodIdleUnits = periodIdleUnits * quotaPercent_ / periodUsedUnits;
    lastUsedUnits_ = context_->usedUnits;
}

void NpuTimesliceScheduler::SchedulerThread(bool &terminating)
{
    while (IsValid()) {
        if (!loaded_) {
            std::this_thread::yield();
            continue;
        }
        if (terminating) {
            return;
        }
        quota_ = TimeSlice();
        currentSlice_ = quota_ * quotaPercent_;
        while (!terminating) {
            if (!CheckCurrent()) {
                break;
            }
            std::this_thread::yield();
        }
#ifdef UNIT_TEST
        if (periodBreak) {
            break;
        }
#endif
        continue;
        auto now = ExecuteTimeslice(begin);
        auto overTime = now - begin - currentSlice_;
        ExecuteIdleTime(quota_);
#ifdef UNIT_TEST
        if (periodBreak) {
            break;
        }
#endif
    }
}