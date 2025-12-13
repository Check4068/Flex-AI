/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2025. All rights reserved.
 */
#ifndef REGISTER_H
#define REGISTER_H

#include <string>

namespace xpu {
bool IsDangerousCommand(const std::string& command)
int ExecCommand(const std::string& command);
bool CheckCgroupData(const std::string& groupData);
int RegisterToDevicePlugin();
void FileOperateErrorHandler(const std::ifstream& file, const std:: string& path);
int GetCgroupData(const std::string& groupPath, std::string& groupData);

#ifdef UNIT_TEST
void SetProcCgroupPath(const std::string& path);
#endif
} // namespace xpu

#endif