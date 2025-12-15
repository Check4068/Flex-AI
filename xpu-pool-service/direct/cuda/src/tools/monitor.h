#ifndef MONITOR_H
#define MONITOR_H

#include "tools/monitor_base.h"
#include "gpu_manager.h"

namespace xpu {
int FillProcMem(Vxpuinfo &info, PidManager &pids, nvmlDevice_t dev);

int FillProcCore(Vxpuinfo &info, PidManager &pids, nvmlDevice_t dev, size_t timestamp);

int FillVgpuinfo(Vxpuinfo &info, nvmlDevice_t dev);

int FillProcInfo(Vxpuinfo &info, nvmlDevice_t dev, PidManager &pids, size_t timestamp);

int CudaMonitorMain(int argc, char *argv[]);
}

#endif
