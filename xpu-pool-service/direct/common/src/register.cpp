/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2025. All rights reserved.
 */
#include <filesystem>
#include <regex>
#include <string>
#include <fstream>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>
#include "log.h"
#include "common.h"
#include "register.h"

using namespace std;
namespace xpu {

#ifndef UNIT_TEST
const static string PROC_CGROUP_PATH = "/proc/self/cgroup";
#else
static string PROC_CGROUP_PATH = "/run/xpu/fake/cgroup";
void SetProcCgroupPath(const string& path)
{
    PROC_CGROUP_PATH = path;
}
#endif

const static string RPC_CLIENT_NAME = "xpu-client-tool";
const static string RPC_CLIENT_PATH = "/opt/xpu/bin/xpu-client-tool";
const static int TRY_TIMES = 10;

void FileOperateErrorHandler(const std::istream& file, const string &path)
{
    if (file.bad()) {
        log_err("I/O error while reading file {%s}", path);
    } else if (!file.eof()) {
        log_err("File {} reached the end", path);
    } else if (!file.fail()) {
        log_err("Non-fatal error occurred while opening {}", path);
    } else {
        log_err("Unexpected error occurred while opening {}", path);
    }
}

int GetCgroupData(const string& groupPath, string& groupData)
{
    // open cgroup file
    ifstream grp(groupPath);
    if (!grp.is_open()) {
        FileOperateErrorHandler(grp, groupPath);
        return RET_FAIL;
    }

    // get memory line
    string memLine;
    const string memoryHeader = "memory:";
    string::size_type pos = memLine.npos;
    while (getline(grp, memLine)) {
        pos = memLine.find(memoryHeader);
        if (pos != memLine.npos) {
            break;
        }
    }
    if (pos == memLine.npos) {
        log_err("cgroup failed");
        return RET_FAIL;
    }

    // get cgroup data
    groupData = memLine.substr(pos + memoryHeader.size());
    if (!CheckCgroupData(groupData)) {
        return RET_FAIL;
    }
    return RET_SUCC;
}

int RegisterWithData(const string& cgroupData)
{
    pid_t pid = fork();
    if (pid < 0) {
        log_err("fork child process failed, errno is {}", strerror(errno));
        return RET_FAIL;
    } else if (pid == 0) {
        // child
        if (IsDangerousCommand(cgroupData)) {
            exit(EXIT_FAILURE);
        }

        if (!std::filesystem::exists(RPC_CLIENT_PATH)) {
            log_err("{} no exist", RPC_CLIENT_PATH);
            exit(EXIT_FAILURE);
        }

        log_info("run: {} --cgroup-path {}", RPC_CLIENT_PATH, cgroupData);
        execl(RPC_CLIENT_PATH.c_str(), RPC_CLIENT_NAME.c_str(),
            "--cgroup-path", cgroupData.data(), nullptr);
        log_err("run rpc client failed, errno is {}", strerror(errno));
        exit(EXIT_FAILURE);
    } else {
        // parent
        int wstatus = 0;
        int wret = waitpid(pid, &wstatus, WUNTRACED | WCONTINUED);
        if (wret == -1) {
            log_err("waitpid failed, error {}", strerror(errno));
            return RET_FAIL;
        }
        if (!WIFEXITED(wstatus) || WEXITSTATUS(wstatus) != 0) {
            log_warn("unexpected exit status {}", wstatus);
            return RET_FAIL;
        }
        log_info("rpc client exit success");
    }

    return RET_SUCC;
}

/*
* (1) Command should not include dangerous command;
* (2) Dangerous command includes: |, &, ;, <, >, /, \, `, \n, \t, *, ?, ", ', (, )
*/
bool IsDangerousCommand(const string& command)
{
    const string blacklist = "|&;<>/\\`\n\t*?\"'()";
    string::size_type pos = command.find_first_of(blacklist);
    if (pos != string::npos) {
        log_err("{%s} is dangerous", command);
        return true;
    }
    log_info("{%s} is safe", command);
    return false;
}

/*
* (1) For systemd, cgroup pattern should be like:
*     11:memory:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-podxxxx.slice/docker-xxxx.scope
*     11:memory:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-podxxxx.slice/cri-containerd-xxxx.scope
*     11:memory:/kubepods.slice/kubepods-besteffort.slice/docker-xxxx.scope
*     11:memory:/kubepods.slice/kubepods-besteffort.slice/cri-containerd-xxxx.scope
* (2) For groups, cgroup pattern should be:
*     11:memory:/kubepods/besteffort/podxxxx/xxxx
*     11:memory:/kubepods/besteffort/podxxxx/xxx
*     11:memory:/contianer.slice/kubepods/besteffort/podxxx/xxx
* (3) The char '.' matches any one char in regex, use "\\." to translate it in c++.
*/
bool CheckCgroupData(const string &groupData)
{
    const string podIdReg = "kubepods-p([0-9a-f]{8})\\.slice";
    const string containerId = "/(cri-containerd|docker)-([0-9a-f]{64})\\.scope";
    regex patternSystemdQos("^/(containerd\\.slice/)?kubepods\\.slice/" + podIdReg + containerId + "$");
    regex patternSystemdBasic("^/(containerd\\.slice/)?kubepods\\.slice/" + containerId + "$");
    regex patternGroupsQos("^/kubepods/([a-z0-9]+)/pod([a-f0-9]{8})/([a-f0-9]{64})$");
    regex patternGroupsBasic("^/kubepods/([a-z0-9]+)/([a-f0-9]{8})/([a-f0-9]{64})$");

    if (regex_match(groupData, patternSystemdQos)) {
        log_info("check qos format success: {%s}", groupData);
        return true;
    }
    if (regex_match(groupData, patternSystemdBasic)) {
        log_info("check basic format success: {%s}", groupData);
        return true;
    }
    if (regex_match(groupData, patternGroupsQos)) {
        log_info("check cgroups qos format success: {%s}", groupData);
        return true;
    }
    if (regex_match(groupData, patternGroupsBasic)) {
        log_info("check cgroups basic format success: {%s}", groupData);
        return true;
    }
    log_err("check format failed: {%s}", groupData);
    return false;
}

int RegisterToDevicePlugin(void)
{
    string groupData;
    int ret = GetCgroupData(PROC_CGROUP_PATH, groupData);
    if (ret != RET_SUCC) {
        log_err("get cgroup data failed, ret is {%d}", ret);
        return RET_FAIL;
    }

    for (int i = 0; i < TRY_TIMES; i++) {
        ret = RegisterWithData(groupData);
        if (ret == RET_SUCC) {
            log_info("register with data success");
            return RET_SUCC;
        }
#ifndef UNIT_TEST
        log_info("register with data failed, retry {%d} time", i + 1);
        this_thread::sleep_for(std::chrono::seconds(1));
#else
        break;
#endif
    }
    return RET_FAIL;
}
} // namespace xpu