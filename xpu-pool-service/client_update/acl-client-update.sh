#!/bin/bash
# ACL 客户端更新脚本
# 功能：初始化 ACL 库文件替换，并持续监控和更新相关文件
# 说明：用于在容器或系统中安装和更新 NPU 监控工具和 ACL 库包装器
# 与 CUDA 版本类似，但针对华为 Ascend NPU 的 ACL（Ascend Computing Language）库

# 定义根路径和工作路径
root_path="/root"                                              # 文件源路径，存放更新文件的根目录
work_path="/opt/xpu"                                           # 工作路径，XPU 工具的安装目录
work_lib_path=${work_path}/lib                                 # 库文件安装路径
work_bin_path=${work_path}/bin                                 # 可执行文件安装路径
lib_path="/usr/local/Ascend/ascend-toolkit/latest/lib64"      # Ascend NPU 工具包的库路径
monitor_name="npu-monitor"                                     # NPU 监控工具文件名
monitor_link="xpu-monitor"                                     # 监控工具的符号链接名称
tool_name="xpu-client-tool"                                    # 客户端工具文件名

# 创建工作目录
mkdir -m 555 -p ${work_lib_path}
mkdir -m 555 -p ${work_bin_path}

# 备份监测列表
declare -A backup_list_all_rx
declare -A backup_list_user_r

# 定义文件路径变量
monitor_path=${work_bin_path}/${monitor_name}        # 监控工具安装路径
monitor_linkpath=${work_bin_path}/${monitor_link}    # 监控工具符号链接路径
tool_path=${work_bin_path}/${tool_name}              # 客户端工具安装路径
root_monitor_path=${root_path}/${monitor_name}       # 根目录下的监控工具路径
root_tool_path=${root_path}/${tool_name}             # 根目录下的客户端工具路径

# 拷贝文件到 /opt/xpu 目录下
install -m 555 ${root_monitor_path} ${monitor_path}
install -m 555 ${root_tool_path} ${tool_path}
# 创建监控工具的符号链接（从 xpu-monitor 指向 npu-monitor）
ln -fs ${monitor_name} ${monitor_linkpath}
# 将安装的文件路径映射添加到备份列表中（用于后续监控更新）
backup_list_all_rx[${monitor_path}]=${root_monitor_path}
backup_list_all_rx[${tool_path}]=${root_tool_path}

make_backup() {
    # 获取系统库文件的真实路径（处理符号链接），参数 $1 是库文件名
    local native=$(readlink -m ${lib_path}/$1)  # 例如：/usr/local/Ascend/ascend-toolkit/latest/lib64/libruntime.so
    # 构建原始备份文件路径，参数 $2 是原始备份文件名
    local original=${work_lib_path}/$2  # 例如：/opt/xpu/lib/libruntime_original.so
    # 构建工作目录备份文件路径
    local backup=${work_lib_path}/$1.bak  # 例如：/opt/xpu/lib/libruntime.so.bak
    # 构建容器根目录备份文件路径
    local container_backup=${root_path}/$1.bak  # 例如：/root/libruntime.so.bak
    # 构建包装器文件路径，参数 $3 是包装器文件名
    local direct=${root_path}/$3  # 例如：/root/libruntime_direct.so

    # 检查系统库文件是否存在
    if [ ! -f "${native}" ]; then
        echo "missing file: ${native}"
        return 1
    fi

    # 如果备份文件不存在，进行首次初始化
    if [ ! -f "${backup}" ]; then
        install -m 555 "${native}" "${original}" \
        && install -m 400 "${native}" "${backup}" \
        && install -m 400 "${backup}" "${container_backup}" \
        && install -m 555 "${direct}" "${native}"
    else
        # 已初始化过：更新容器内的备份文件
        install -m 400 "${backup}" "${container_backup}"
    fi

    # 将文件路径映射添加到备份列表中，用于后续监控和更新
    backup_list_all_rx["${native}"]="${original}"      # 系统库 -> 原始备份
    backup_list_all_rx["${original}"]="${backup}"      # 原始备份 -> 工作目录备份
    backup_list_user_r["${backup}"]="${container_backup}"  # 工作目录备份 -> 容器根目录备份
}

# 调用 make_backup 函数，备份并替换 libruntime.so 库文件
# 参数：库文件名、原始备份文件名、包装器文件名
make_backup libruntime.so libruntime_original.so libruntime_direct.so \
&& echo "client file initialization completed" \
|| (echo "client file initialization failed" && exit 1)

update_all_rx() {
    ...
}

# 主循环：持续监控文件更新
while true; do
    # 更新所有需要监控的文件（包括库文件和工具文件）
    update_all_rx
    # 更新用户级别的备份文件
    update_user_r
    
    # 检查监控工具的符号链接是否正确，如果不正确则恢复
    if [ "$(readlink ${monitor_linkpath})" != "${monitor_path}" ]; then
        echo "monitor link is restored"
        ln -fs ${monitor_name} ${monitor_linkpath}
    fi
    
    # 休眠 5 秒后继续监控
    sleep 5
    echo 'files is being monitored'
done
