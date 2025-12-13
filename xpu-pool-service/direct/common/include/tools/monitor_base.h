/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */
#ifndef MONITOR_BASE_H
#define MONITOR_BASE_H

#include <cctype>
#include <future>
#include <map>
#include "resource_config.h"
#include "log.h"

namespace xpu {

enum class OutputFormat : char {
    NONE = '\0',
    TABLE = 't',
    JSON = 'j'
};

enum class VxpuType : char {
    VGPU = 'G',
    VNPU = 'N',
};

constexpr int PERIOD_DEFAULT = 60;   // one minute
constexpr int PERIOD_MIN = 1;        // one second
constexpr int PERIOD_MAX = 1024 * 60 * 24; // one day
constexpr int MAX_PIDS = 1024;

struct Args {
    int period = PERIOD_DEFAULT;
    OutputFormat format = OutputFormat::TABLE;
};

int ParseArgs(Args& args, int argc, char const* argv[]);

struct OutputFormatter {
    OutputFormat format = OutputFormat::NONE;

    constexpr auto parse(fmt::format_parse_context& ctx) {
        for (auto it = ctx.begin(); it != ctx.end(); ++it) {
            if (*it == 'j') {
                format = OutputFormat::JSON;
            } else if (*it == 't') {
                format = OutputFormat::TABLE;
            } else {
                ctx.on_error("invalid output format");
            }
            return it;
        }
        return ctx.end();
    }
};

struct ProcessInfo {
    uint32_t core = 0;
    size_t memory = 0;
};

struct VxpuInfo {
    VxpuType type;
    uint32_t id = 0;
    uint32_t core = 0;
    uint32_t coreQuota = 0;
    size_t memory = 0;
    size_t memoryQuota = 0;
    std::map<uint32_t, ProcessInfo> processes;

    VxpuInfo(ResourceConfig &config, VxpuType type, int32_t id) : type(type), id(id)
    {
        if (config.LimitComputingPower()) {
            coreQuota = config.ComputingPowerQuota();
        } else {
            coreQuota = PERCENT_MAX;
        }
        if (config.LimitMemory()) {
            memoryQuota = config.MemoryQuota();
        } else {
            memoryQuota = 0;
        }
    }

TESTABLE_PRIVATE:
    VxpuInfo(VxpuType type) : type(type), coreQuota(PERCENT_MAX), memoryQuota(0) {}
};

struct ContainerVxpuInfo {
    VxpuType type;
    std::vector<VxpuInfo> vxpus;

    ContainerVxpuInfo(VxpuType type) : type(type)
    {}
};

// namespace xpu
template <>
class fmt::formatter<std::pair<uint32_t, xpu::ProcessInfo>> : public xpu::VxpuFormatter {
public:
    template <typename Context>
    auto format(const std::pair<uint32_t, xpu::ProcessInfo> &info, Context &ctx) const
    {
        if (format_ == xpu::OutputFormat::JSON) {
            return format_to(ctx.out(),
                "{{\"pid\": {:d}, \"core\": {:d}, \"memory\": {:d}}}",
                info.first,
                info.second.core,
                info.second.memory);
        } else {
            return format_to(ctx.out(),
                "pid: {:d}, core usage: {:d}%, memory usage: {:d}MB",
                info.first,
                info.second.core,
                info.second.memory / MEGABYTE);
        }
    }
};

template <>
class fmt::formatter<xpu::VxpuInfo> : public xpu::VxpuFormatter {
public:
    template <typename Context>
    auto format(const xpu::VxpuInfo &info, Context &ctx) const
    {
        if (format_ == xpu::OutputFormat::JSON) {
            return format_to(ctx.out(),
                "{{\"device\": {:d}, \"id\": {:d}, \"core_quota\": {:d}, \"memory\": {:d}, \"memory_quota\": {:d},\n"
                "\"processes\": [{}]}}",
                info.type,
                info.id,
                info.coreQuota,
                info.memory,
                info.memoryQuota,
                fmt::join(info.processes, ",\n"));
        } else {
            return format_to(ctx.out(),
                "v{}PU {} usage: {:02}%, limit: {:02}%, memory usage: {:6}/{}MB\n{:t}",
                char(info.type)
                info.id,
                char(info.type),
                info.coreQuota,
                info.memory / MEGABYTE,
                fmt::join(info.processes, "\n\t"));
        }
    }
};

template <>
class fmt::formatter<xpu::ContainerVxpuInfo> : public xpu::VxpuFormatter {
public:
    template <typename Context>
    auto format(const xpu::ContainerVxpuInfo &info, Context &ctx) const
    {
        auto fmt = this->OutputFormat();
        if (fmt == xpu::OutputFormat::JSON) {
            return fmt::format_to(ctx.out(),
                "{{\"type\": \"v{}PU\", \"vxpus\": [\n{:j}\n]}}",
                char(info.type),
                fmt::join(info.vxpus, ",\n"));
        } else {
            return fmt::format_to(ctx.out(),
                "v{}PU num: {}\n{:t}", char(info.type), info.vxpus.size(), fmt::join(info.vxpus, "\n")",
        }
    }
};

#endif