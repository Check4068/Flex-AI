# Diff Details

Date : 2025-12-13 16:40:57

Directory c:\\Users\\admin\\IdeaProjects\\untitled\\scheduler\\xpu-pool-service\\GPU-device-plugin\\pkg

Total : 46 files,  2076 codes, 531 comments, 416 blanks, all 3023 lines

[Summary](results.md) / [Details](details.md) / [Diff Summary](diff.md) / Diff Details

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [xpu-pool-service/GPU-device-plugin/install/gpu\_env\_init.sh](/xpu-pool-service/GPU-device-plugin/install/gpu_env_init.sh) | Shell Script | -21 | -4 | -5 | -30 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/Chart.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/Chart.yaml) | YAML | -7 | -1 | -1 | -9 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/gpu-client-update-daemonset.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/gpu-client-update-daemonset.yaml) | YAML | -53 | -1 | 0 | -54 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/gpu-device-plugin-configmap.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/gpu-device-plugin-configmap.yaml) | YAML | -10 | -2 | 0 | -12 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/gpu-device-plugin-daemonset.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/gpu-device-plugin-daemonset.yaml) | YAML | -94 | -1 | 0 | -95 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/npu-client-update-daemonset.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/npu-client-update-daemonset.yaml) | YAML | -53 | -1 | 0 | -54 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/npu-device-plugin-daemonset.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/npu-device-plugin-daemonset.yaml) | YAML | -86 | -1 | 0 | -87 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/role-binding.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/role-binding.yaml) | YAML | -15 | -1 | 0 | -16 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/role.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/role.yaml) | YAML | -33 | -1 | -1 | -35 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/runtimeclass.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/runtimeclass.yaml) | YAML | -5 | 0 | 0 | -5 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/service-account.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/service-account.yaml) | YAML | -7 | -1 | -1 | -9 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/xpu-exporter-daemonset.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/xpu-exporter-daemonset.yaml) | YAML | -67 | -1 | 0 | -68 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/xpu-exporter-service.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/xpu-exporter-service.yaml) | YAML | -15 | -1 | 0 | -16 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/xpu-servicemonitor.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/xpu-servicemonitor.yaml) | YAML | -22 | -1 | 0 | -23 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/values.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/values.yaml) | YAML | -97 | -1 | -13 | -111 |
| [xpu-pool-service/GPU-device-plugin/install/install.sh](/xpu-pool-service/GPU-device-plugin/install/install.sh) | Shell Script | -105 | -12 | -17 | -134 |
| [xpu-pool-service/GPU-device-plugin/install/npu\_env\_init.sh](/xpu-pool-service/GPU-device-plugin/install/npu_env_init.sh) | Shell Script | -3 | -2 | -2 | -7 |
| [xpu-pool-service/GPU-device-plugin/install/uninstall.sh](/xpu-pool-service/GPU-device-plugin/install/uninstall.sh) | Shell Script | -48 | -7 | -10 | -65 |
| [xpu-pool-service/GPU-device-plugin/install/yaml/client-update.yaml](/xpu-pool-service/GPU-device-plugin/install/yaml/client-update.yaml) | YAML | -37 | -1 | -1 | -39 |
| [xpu-pool-service/GPU-device-plugin/install/yaml/gpu-device-plugin.yaml](/xpu-pool-service/GPU-device-plugin/install/yaml/gpu-device-plugin.yaml) | YAML | -118 | -4 | -9 | -131 |
| [xpu-pool-service/GPU-device-plugin/install/yaml/namespace.yaml](/xpu-pool-service/GPU-device-plugin/install/yaml/namespace.yaml) | YAML | -6 | -1 | -1 | -8 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/cgo\_helpers\_static.go](/xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/cgo_helpers_static.go) | Go | 12 | 3 | 3 | 18 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml\_api.go](/xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml_api.go) | Go | 41 | 22 | 19 | 82 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml\_common.go](/xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml_common.go) | Go | 0 | 0 | 1 | 1 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml\_const.go](/xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml_const.go) | Go | 48 | 15 | 13 | 76 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml\_const\_static.go](/xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml_const_static.go) | Go | 18 | 5 | 6 | 29 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml\_device.go](/xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml_device.go) | Go | 110 | 4 | 18 | 132 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml\_event.go](/xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml_event.go) | Go | 43 | 7 | 8 | 58 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml\_lib.go](/xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml_lib.go) | Go | 110 | 4 | 27 | 141 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml\_refcount.go](/xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml_refcount.go) | Go | 12 | 7 | 5 | 24 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml\_wrap.go](/xpu-pool-service/GPU-device-plugin/pkg/api/gonvml/nvml_wrap.go) | Go | 169 | 244 | 39 | 452 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/graph/graph.go](/xpu-pool-service/GPU-device-plugin/pkg/api/graph/graph.go) | Go | 33 | 17 | 7 | 57 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/lock/nodelock.go](/xpu-pool-service/GPU-device-plugin/pkg/api/lock/nodelock.go) | Go | 164 | 8 | 11 | 183 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/log/console.go](/xpu-pool-service/GPU-device-plugin/pkg/api/log/console.go) | Go | 40 | 7 | 10 | 57 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/log/file.go](/xpu-pool-service/GPU-device-plugin/pkg/api/log/file.go) | Go | 361 | 64 | 62 | 487 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/log/logger.go](/xpu-pool-service/GPU-device-plugin/pkg/api/log/logger.go) | Go | 153 | 23 | 31 | 207 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/api.pb.go](/xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/api.pb.go) | Go | 289 | 18 | 42 | 349 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/api\_grpc.pb.go](/xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/api_grpc.pb.go) | Go | 108 | 22 | 18 | 148 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/service\_impl.go](/xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/service_impl.go) | Go | 422 | 21 | 33 | 476 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/service\_impl\_test.go](/xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/service_impl_test.go) | Go | 233 | 4 | 22 | 259 |
| [xpu-pool-service/GPU-device-plugin/pkg/plugin/cache.go](/xpu-pool-service/GPU-device-plugin/pkg/plugin/cache.go) | Go | 62 | 11 | 14 | 87 |
| [xpu-pool-service/GPU-device-plugin/pkg/plugin/config/config.go](/xpu-pool-service/GPU-device-plugin/pkg/plugin/config/config.go) | Go | 8 | 9 | 3 | 20 |
| [xpu-pool-service/GPU-device-plugin/pkg/plugin/plugin.go](/xpu-pool-service/GPU-device-plugin/pkg/plugin/plugin.go) | Go | 343 | 29 | 49 | 421 |
| [xpu-pool-service/GPU-device-plugin/pkg/plugin/register.go](/xpu-pool-service/GPU-device-plugin/pkg/plugin/register.go) | Go | 118 | 7 | 17 | 142 |
| [xpu-pool-service/GPU-device-plugin/pkg/plugin/types/types.go](/xpu-pool-service/GPU-device-plugin/pkg/plugin/types/types.go) | Go | 68 | 19 | 14 | 101 |
| [xpu-pool-service/GPU-device-plugin/pkg/plugin/xpu/xpu\_common.go](/xpu-pool-service/GPU-device-plugin/pkg/plugin/xpu/xpu_common.go) | Go | 13 | 6 | 5 | 24 |

[Summary](results.md) / [Details](details.md) / [Diff Summary](diff.md) / Diff Details