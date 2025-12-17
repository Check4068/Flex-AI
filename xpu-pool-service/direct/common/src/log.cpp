/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */
#include <string>
#include <vector>
#include <filesystem>
#include <iostream>
#include <regex>
#include <sys/stat.h>
#include <spdlog/spdlog.h>
#include <spdlog/sinks/stdout_color_sinks.h>
#include <spdlog/sinks/rotating_file_sink.h>
#include <spdlog/cfg/env.h>
#include "common.h"
#include "log.h"
#include "register.h"

namespace fs = std::filesystem;

namespace xpu {
std::string GetContainerIDFromCgroup(const std::string &filePath)
{
    std::regex containerIdPattern("[0-9a-f]{64}");

    std::string cgroupData;
    int ret = GetCgroupData(filePath, cgroupData);
    if (ret != RET_SUCC) {
        return "";
    }

    std::smatch match;
    if (std::regex_search(cgroupData, match, containerIdPattern)) {
        return match[0].str();
    } else {
        std::cerr << "Get container id failed" << std::endl;
    }

    return "";
}

void LogToFile(const std::string &logdir, std::shared_ptr<spdlog::logger> &logger)
{ 
    std::string FilePath = "/proc/self/cgroup";
    const int CNTR_ID_CUT_LEN = 8;
    std::string containerId = GetContainerIDFromCgroup(FilePath).empty() ? "nocontainer" : GetContainerIDFromCgroup(FilePath).substr(0, CNTR_ID_CUT_LEN);
    pid_t pid = getpid();
    std::string fileName = logdir + containerId + "-" + std::to_string(pid) + ".log";
    log->sinks().push_back(std::make_shared<spdlog::sinks::rotating_file_sink_mt>(fileName, 1024 * 1024 * 5, 10));
}

void logInit(const std::string &loggerName, const std::string &sourceId)
{ 
    setenv("SPDLOG_LEVEL", SPDLOG_LEVEL_NAME_WARNING.data(), 0);
    std::shared_ptr<spdlog::logger> xpuLogger;
    xpuLogger = std::make_shared<spdlog::logger>(loggerName + "-" + sourceId);
    xpuLogger->sinks().push_back(std::make_shared<spdlog::sinks::stdout_color_sink_mt>());
    const std::string logdir = "/var/log/xpu/";
    if (fs::exists(logdir) && fs::is_directory(logdir)) {
        LogToFile(logdir, xpuLogger);
    } 
    xpuLogger->flush_on(spdlog::level::warn);
    spdlog::register_logger(xpuLogger);
    spdlog::set_default_logger(xpuLogger);
    spdlog::cfg::load_env_levels();
}
}