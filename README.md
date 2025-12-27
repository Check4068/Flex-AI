# Flex::AI  


## Flex::AI 编译

1. 安装依赖项

下载 [https://github.com/gabime/spdlog/archive/refs/tags/v1.12.0.zip](https://github.com/gabime/spdlog/archive/refs/tags/v1.12.0.zip) 到本地，执行：

```Bash

```

下载 [https://github.com/NixOS/patchelf/archive/refs/tags/0.18.0.zip](https://github.com/NixOS/patchelf/archive/refs/tags/0.18.0.zip) 到本地，执行：

```Bash

```

1. 编译

下载 [https://github.com/Check4068/Flex-AI/archive/refs/heads/main.zip](https://github.com/Check4068/Flex-AI/archive/refs/heads/main.zip) 到本地，执行：

```Bash

```

如果一切顺利，执行产物在`direct/cuda`目录下，为`libcuda_direct.so`。

## 设备插件

go的版本为1.25.0，建议保持一致：

```Bash

```

## 调度组件

下载 [https://github.com/volcano-sh/volcano/archive/refs/tags/v1.10.2.zip](https://github.com/volcano-sh/volcano/archive/refs/tags/v1.10.2.zip) 到本地，执行：

```Bash

```

将解压后的`volcano-1.10.2.zip`复制到创建出的目录下：

```Bash

```

将`xpu-scheduler-plugin`目录整个复制到volcano的plugins目录下：

```Bash

```

执行编译：

```Bash

```

在`out`目录下可见编译产物为`huawei-xpu.so`、`vc-scheduler`、`vc-controller-manager`。

## GPU虚拟化源码编译

### 前置准备

1. 安装cuda，请使用下述命令：

```Bash

```

1. 安装相关依赖项

2.1 下载 [https://github.com/gabime/spdlog/archive/refs/tags/v1.12.0.zip](https://github.com/gabime/spdlog/archive/refs/tags/v1.12.0.zip) 到本地，执行下面命令：

```Bash

```

2.2 下载 [https://github.com/NixOS/patchelf/archive/refs/tags/0.18.0.zip](https://github.com/NixOS/patchelf/archive/refs/tags/0.18.0.zip) 到本地，执行下面命令：

```Bash

```

### 3、编译

下载 [https://github.com/Check4068/Flex-AI/archive/refs/heads/main.zip](https://github.com/Check4068/Flex-AI/archive/refs/heads/main.zip) 到本地，执行下面命令：

```Bash

```

上述步骤顺利执行后，执行产物在`direct/cuda`目录下，为`libcuda_direct.so`。

## GPU虚拟化试用（正在完善中）

1. 打包`libcuda_direct.so`到镜像

将构建产物`libcuda_direct.so`和`cuda-client-update.sh`放在同一个目录下，编写Dockerfile，示例如下：

```Dockerfile

```

脚本的地址为：

[Flex-AI/xpu-pool-service/client_update/cuda-client-update.sh at main · Check4068/Flex-AI · GitHub](https://github.com/Check4068/Flex-AI/blob/main/xpu-pool-service/client_update/cuda-client-update.sh)

在Dockerfile所在的目录执行：

```Bash

```

1. GPU-device-plugin的编译容器

在GPU-device-plugin的根目录（Makefile所在的目录）执行`make -j`，编译产物为`gpu-device-plugin`。

打包流程与上述类似，将构建产物与Dockerfile放在同一目录下，Dockerfile示例如下：

```Dockerfile

```

在Dockerfile所在的目录执行：

```Bash

```

1. yaml文件编写

将运行时的镜像注入到要进行负载的gpu节点上，编写yaml/helm chart 推进调度任务（上述业务的yaml文件中得有gpu资源，示例如下：

```YAML

```

1. 使用命令查看算力利用率

```Bash

```

查看算力利用率是否在20%上下波动

**备注**：vggu最多支持切分5份，申请一个vggu表示使用20%算力

## 部署流程

1. 上传模型权重到希望调度的有xpu卡的工作节点的任意目录下

模型文件：`DeepSeek-R1-Distill-Llama-8B`

1. 拉取vllm镜像

```Bash

```

1. 登陆master节点，给有xpu卡的工作节点上标签

```Bash

```

1. 在master节点创建命名空间

```Bash

```

1. 在master节点编辑并上传`deployment.yaml`和`service.yaml`

    - **deployment.yaml**

```YAML

```

- **service.yaml**

```YAML

```

## 启动

将镜像文件放入对应的本地挂载路径之后，执行：

```Bash

```

检查pod状态：

```Bash

```

容器创建成功之后，查看执行日志：

```Bash

```

出现类似于：`INFO 03-27 23:19:15 api_server.py:958] Starting vLLM API server on http://0.0.0.0:8000`，即为成功。

查看GPU卡使用情况：

```Bash

```

## 调用

首先获取IP地址：

```Bash

```

获取接口说明：

```Bash

```

聊天：

```Bash

```

- 如果配置了`hostNetwork: true`也可以用postman请求。

## 补充：vllm参数说明

|参数|说明|默认值|
|---|---|---|
|load_format|模型权重加载的格式|"auto"|
|gpu-memory-utilization|用于模型执行器的GPU内存分配，范围0到1|0.9|
|max-model-len|模型上下文长度，如果未指定，将自动从模型配置中推导|config.json#max_position_embeddings|
|tensor-parallel-size|张量并行的副本数量|-|
- vLLM默认以通16GB，剩余空间主要用于KV缓存。

- `gpu-memory-utilization`参数下调后，如果出现OOM，需要根据错误信息对应下调`--max-model-len`参数

- 全分仓 → 全部预分配 → 分配点0

```Plain Text


```

