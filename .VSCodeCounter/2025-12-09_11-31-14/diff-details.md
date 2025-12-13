# Diff Details

Date : 2025-12-09 11:31:14

Directory c:\\Users\\admin\\IdeaProjects\\untitled\\scheduler\\xpu-pool-service\\GPU-device-plugin\\install

Total : 89 files,  -4447 codes, -451 comments, -635 blanks, all -5533 lines

[Summary](results.md) / [Details](details.md) / [Diff Summary](diff.md) / Diff Details

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [xpu-pool-service/CMakeLists.txt](/xpu-pool-service/CMakeLists.txt) | CMake | -19 | 0 | -7 | -26 |
| [xpu-pool-service/GPU-device-plugin/Makefile](/xpu-pool-service/GPU-device-plugin/Makefile) | Makefile | -34 | -2 | -6 | -42 |
| [xpu-pool-service/GPU-device-plugin/client/client.go](/xpu-pool-service/GPU-device-plugin/client/client.go) | Go | -53 | -13 | -9 | -75 |
| [xpu-pool-service/GPU-device-plugin/client/client\_test.go](/xpu-pool-service/GPU-device-plugin/client/client_test.go) | Go | -146 | -39 | -17 | -202 |
| [xpu-pool-service/GPU-device-plugin/cmd/main.go](/xpu-pool-service/GPU-device-plugin/cmd/main.go) | Go | -97 | -38 | -25 | -160 |
| [xpu-pool-service/GPU-device-plugin/cmd/main\_test.go](/xpu-pool-service/GPU-device-plugin/cmd/main_test.go) | Go | -347 | -2 | -9 | -358 |
| [xpu-pool-service/GPU-device-plugin/examples/vgpu-exceed.yml](/xpu-pool-service/GPU-device-plugin/examples/vgpu-exceed.yml) | YAML | -14 | -1 | 0 | -15 |
| [xpu-pool-service/GPU-device-plugin/examples/vgpu-test.yml](/xpu-pool-service/GPU-device-plugin/examples/vgpu-test.yml) | YAML | -25 | 0 | 0 | -25 |
| [xpu-pool-service/GPU-device-plugin/examples/vgpu-whole-card.yml](/xpu-pool-service/GPU-device-plugin/examples/vgpu-whole-card.yml) | YAML | -13 | 0 | 0 | -13 |
| [xpu-pool-service/GPU-device-plugin/go.mod](/xpu-pool-service/GPU-device-plugin/go.mod) | Go Module File | -63 | 0 | -4 | -67 |
| [xpu-pool-service/GPU-device-plugin/go.sum](/xpu-pool-service/GPU-device-plugin/go.sum) | Go Checksum File | -177 | 0 | 0 | -177 |
| [xpu-pool-service/GPU-device-plugin/install/gpu\_env\_init.sh](/xpu-pool-service/GPU-device-plugin/install/gpu_env_init.sh) | Shell Script | 21 | 4 | 5 | 30 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/Chart.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/Chart.yaml) | YAML | 7 | 1 | 1 | 9 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/gpu-client-update-daemonset.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/gpu-client-update-daemonset.yaml) | YAML | 53 | 1 | 0 | 54 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/gpu-device-plugin-configmap.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/gpu-device-plugin-configmap.yaml) | YAML | 10 | 2 | 0 | 12 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/gpu-device-plugin-daemonset.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/gpu-device-plugin-daemonset.yaml) | YAML | 94 | 1 | 0 | 95 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/npu-client-update-daemonset.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/npu-client-update-daemonset.yaml) | YAML | 53 | 1 | 0 | 54 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/npu-device-plugin-daemonset.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/npu-device-plugin-daemonset.yaml) | YAML | 86 | 1 | 0 | 87 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/role-binding.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/role-binding.yaml) | YAML | 15 | 1 | 0 | 16 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/role.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/role.yaml) | YAML | 33 | 1 | 1 | 35 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/runtimeclass.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/runtimeclass.yaml) | YAML | 5 | 0 | 0 | 5 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/service-account.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/service-account.yaml) | YAML | 7 | 1 | 1 | 9 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/xpu-exporter-daemonset.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/xpu-exporter-daemonset.yaml) | YAML | 67 | 1 | 0 | 68 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/xpu-exporter-service.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/xpu-exporter-service.yaml) | YAML | 15 | 1 | 0 | 16 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/xpu-servicemonitor.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/templates/xpu-servicemonitor.yaml) | YAML | 22 | 1 | 0 | 23 |
| [xpu-pool-service/GPU-device-plugin/install/helm/gpupool/values.yaml](/xpu-pool-service/GPU-device-plugin/install/helm/gpupool/values.yaml) | YAML | 97 | 1 | 13 | 111 |
| [xpu-pool-service/GPU-device-plugin/install/install.sh](/xpu-pool-service/GPU-device-plugin/install/install.sh) | Shell Script | 105 | 12 | 17 | 134 |
| [xpu-pool-service/GPU-device-plugin/install/npu\_env\_init.sh](/xpu-pool-service/GPU-device-plugin/install/npu_env_init.sh) | Shell Script | 3 | 2 | 2 | 7 |
| [xpu-pool-service/GPU-device-plugin/install/uninstall.sh](/xpu-pool-service/GPU-device-plugin/install/uninstall.sh) | Shell Script | 48 | 7 | 10 | 65 |
| [xpu-pool-service/GPU-device-plugin/install/yaml/client-update.yaml](/xpu-pool-service/GPU-device-plugin/install/yaml/client-update.yaml) | YAML | 37 | 1 | 1 | 39 |
| [xpu-pool-service/GPU-device-plugin/install/yaml/gpu-device-plugin.yaml](/xpu-pool-service/GPU-device-plugin/install/yaml/gpu-device-plugin.yaml) | YAML | 118 | 4 | 9 | 131 |
| [xpu-pool-service/GPU-device-plugin/install/yaml/namespace.yaml](/xpu-pool-service/GPU-device-plugin/install/yaml/namespace.yaml) | YAML | 6 | 1 | 1 | 8 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/api.pb.go](/xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/api.pb.go) | Go | -289 | -18 | -42 | -349 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/api\_grpc.pb.go](/xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/api_grpc.pb.go) | Go | -108 | -22 | -18 | -148 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/service\_impl.go](/xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/service_impl.go) | Go | -422 | -21 | -33 | -476 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/service\_impl\_test.go](/xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/service_impl_test.go) | Go | -233 | -4 | -22 | -259 |
| [xpu-pool-service/GPU-device-plugin/watchers/watchers.go](/xpu-pool-service/GPU-device-plugin/watchers/watchers.go) | Go | -24 | -6 | -7 | -37 |
| [xpu-pool-service/GPU-device-plugin/watchers/watchers\_test.go](/xpu-pool-service/GPU-device-plugin/watchers/watchers_test.go) | Go | -54 | -4 | -10 | -68 |
| [xpu-pool-service/README.md](/xpu-pool-service/README.md) | Markdown | 0 | 0 | -1 | -1 |
| [xpu-pool-service/ci/VersionSet.xml](/xpu-pool-service/ci/VersionSet.xml) | XML | -1 | 0 | 0 | -1 |
| [xpu-pool-service/ci/app\_define.json](/xpu-pool-service/ci/app_define.json) | JSON | -13 | 0 | 0 | -13 |
| [xpu-pool-service/ci/at/at\_deploy.sh](/xpu-pool-service/ci/at/at_deploy.sh) | Shell Script | -23 | -3 | -5 | -31 |
| [xpu-pool-service/ci/at/at\_deploy.yml](/xpu-pool-service/ci/at/at_deploy.yml) | YAML | -48 | 0 | -5 | -53 |
| [xpu-pool-service/ci/build.sh](/xpu-pool-service/ci/build.sh) | Shell Script | -80 | -4 | -18 | -102 |
| [xpu-pool-service/ci/build.yml](/xpu-pool-service/ci/build.yml) | YAML | -44 | 0 | -4 | -48 |
| [xpu-pool-service/ci/buildinfo.sh](/xpu-pool-service/ci/buildinfo.sh) | Shell Script | -11 | -3 | -2 | -16 |
| [xpu-pool-service/ci/cmc/openSource\_x86.xml](/xpu-pool-service/ci/cmc/openSource_x86.xml) | XML | -24 | 0 | 0 | -24 |
| [xpu-pool-service/ci/cmc/upload\_cmc.xml](/xpu-pool-service/ci/cmc/upload_cmc.xml) | XML | -24 | 0 | 0 | -24 |
| [xpu-pool-service/ci/cms\_signature.sh](/xpu-pool-service/ci/cms_signature.sh) | Shell Script | -55 | -3 | -6 | -64 |
| [xpu-pool-service/ci/dependency.xml](/xpu-pool-service/ci/dependency.xml) | XML | -7 | -1 | 0 | -8 |
| [xpu-pool-service/ci/hwp7s\_signature.sh](/xpu-pool-service/ci/hwp7s_signature.sh) | Shell Script | -37 | -4 | -3 | -44 |
| [xpu-pool-service/ci/opensource.xml](/xpu-pool-service/ci/opensource.xml) | XML | -6 | 0 | 0 | -6 |
| [xpu-pool-service/ci/xpu\_pool/build\_x86.yml](/xpu-pool-service/ci/xpu_pool/build_x86.yml) | YAML | -65 | 0 | -5 | -70 |
| [xpu-pool-service/ci/xpu\_pool/build\_xpu\_package.sh](/xpu-pool-service/ci/xpu_pool/build_xpu_package.sh) | Shell Script | -107 | -3 | -13 | -123 |
| [xpu-pool-service/ci/xpu\_pool/xpu\_docker\_build/acl\_client/Dockerfile](/xpu-pool-service/ci/xpu_pool/xpu_docker_build/acl_client/Dockerfile) | Docker | -5 | -2 | -3 | -10 |
| [xpu-pool-service/ci/xpu\_pool/xpu\_docker\_build/cuda\_client/Dockerfile](/xpu-pool-service/ci/xpu_pool/xpu_docker_build/cuda_client/Dockerfile) | Docker | -5 | -2 | -3 | -10 |
| [xpu-pool-service/ci/xpu\_pool/xpu\_docker\_build/exporter/Dockerfile](/xpu-pool-service/ci/xpu_pool/xpu_docker_build/exporter/Dockerfile) | Docker | -20 | -4 | -6 | -30 |
| [xpu-pool-service/ci/xpu\_pool/xpu\_docker\_build/gpu\_device\_plugin/Dockerfile](/xpu-pool-service/ci/xpu_pool/xpu_docker_build/gpu_device_plugin/Dockerfile) | Docker | -10 | -4 | -5 | -19 |
| [xpu-pool-service/ci/xpu\_pool/xpu\_docker\_build/npu\_device\_plugin/Dockerfile](/xpu-pool-service/ci/xpu_pool/xpu_docker_build/npu_device_plugin/Dockerfile) | Docker | -5 | -10 | -4 | -19 |
| [xpu-pool-service/client\_update/acl-client-update.sh](/xpu-pool-service/client_update/acl-client-update.sh) | Shell Script | -60 | -28 | -15 | -103 |
| [xpu-pool-service/client\_update/cuda-client-update.sh](/xpu-pool-service/client_update/cuda-client-update.sh) | Shell Script | -68 | -35 | -12 | -115 |
| [xpu-pool-service/direct/CMakeLists.txt](/xpu-pool-service/direct/CMakeLists.txt) | CMake | -21 | 0 | -4 | -25 |
| [xpu-pool-service/direct/acl/CMakeLists.txt](/xpu-pool-service/direct/acl/CMakeLists.txt) | CMake | -33 | 0 | -7 | -40 |
| [xpu-pool-service/direct/make\_lib\_original.sh](/xpu-pool-service/direct/make_lib_original.sh) | Shell Script | -33 | -1 | -7 | -41 |
| [xpu-pool-service/xpu-exporter/Makefile](/xpu-pool-service/xpu-exporter/Makefile) | Makefile | -20 | -2 | -5 | -27 |
| [xpu-pool-service/xpu-exporter/cmd/xpu-exporter/main.go](/xpu-pool-service/xpu-exporter/cmd/xpu-exporter/main.go) | Go | -107 | -15 | -18 | -140 |
| [xpu-pool-service/xpu-exporter/collector/collector\_service.go](/xpu-pool-service/xpu-exporter/collector/collector_service.go) | Go | -11 | -8 | -6 | -25 |
| [xpu-pool-service/xpu-exporter/collector/gpuservice/gpu\_collector.go](/xpu-pool-service/xpu-exporter/collector/gpuservice/gpu_collector.go) | Go | -200 | -6 | -22 | -228 |
| [xpu-pool-service/xpu-exporter/collector/gpuservice/gpu\_collector\_service.go](/xpu-pool-service/xpu-exporter/collector/gpuservice/gpu_collector_service.go) | Go | -73 | -9 | -11 | -93 |
| [xpu-pool-service/xpu-exporter/collector/npuservice/npu\_collector.go](/xpu-pool-service/xpu-exporter/collector/npuservice/npu_collector.go) | Go | -32 | -6 | -8 | -46 |
| [xpu-pool-service/xpu-exporter/collector/npuservice/npu\_collector\_service.go](/xpu-pool-service/xpu-exporter/collector/npuservice/npu_collector_service.go) | Go | -31 | -9 | -10 | -50 |
| [xpu-pool-service/xpu-exporter/common/cache/lrucache.go](/xpu-pool-service/xpu-exporter/common/cache/lrucache.go) | Go | -319 | -21 | -26 | -366 |
| [xpu-pool-service/xpu-exporter/common/cache/lrucache\_test.go](/xpu-pool-service/xpu-exporter/common/cache/lrucache_test.go) | Go | -71 | -6 | -12 | -89 |
| [xpu-pool-service/xpu-exporter/common/client/client.go](/xpu-pool-service/xpu-exporter/common/client/client.go) | Go | -36 | -5 | -6 | -47 |
| [xpu-pool-service/xpu-exporter/common/client/client\_test.go](/xpu-pool-service/xpu-exporter/common/client/client_test.go) | Go | -42 | -4 | -11 | -57 |
| [xpu-pool-service/xpu-exporter/common/limiter/limit\_handler.go](/xpu-pool-service/xpu-exporter/common/limiter/limit_handler.go) | Go | -223 | -20 | -28 | -271 |
| [xpu-pool-service/xpu-exporter/common/limiter/limit\_handler\_test.go](/xpu-pool-service/xpu-exporter/common/limiter/limit_handler_test.go) | Go | -90 | -6 | -12 | -108 |
| [xpu-pool-service/xpu-exporter/common/limiter/limit\_listener.go](/xpu-pool-service/xpu-exporter/common/limiter/limit_listener.go) | Go | -112 | -8 | -17 | -137 |
| [xpu-pool-service/xpu-exporter/common/limiter/limit\_listener\_test.go](/xpu-pool-service/xpu-exporter/common/limiter/limit_listener_test.go) | Go | -78 | 0 | -16 | -94 |
| [xpu-pool-service/xpu-exporter/common/service/api.pb.go](/xpu-pool-service/xpu-exporter/common/service/api.pb.go) | Go | -292 | -15 | -43 | -350 |
| [xpu-pool-service/xpu-exporter/common/service/api\_grpc.pb.go](/xpu-pool-service/xpu-exporter/common/service/api_grpc.pb.go) | Go | -108 | -24 | -20 | -152 |
| [xpu-pool-service/xpu-exporter/common/test\_utils.go](/xpu-pool-service/xpu-exporter/common/test_utils.go) | Go | -50 | -7 | -10 | -67 |
| [xpu-pool-service/xpu-exporter/common/utils/ip\_utils.go](/xpu-pool-service/xpu-exporter/common/utils/ip_utils.go) | Go | -25 | -7 | -6 | -38 |
| [xpu-pool-service/xpu-exporter/common/utils/type.go](/xpu-pool-service/xpu-exporter/common/utils/type.go) | Go | -31 | -7 | -4 | -42 |
| [xpu-pool-service/xpu-exporter/go.mod](/xpu-pool-service/xpu-exporter/go.mod) | Go Module File | -33 | 0 | -4 | -37 |
| [xpu-pool-service/xpu-exporter/go.sum](/xpu-pool-service/xpu-exporter/go.sum) | Go Checksum File | -42 | 0 | 0 | -42 |
| [xpu-pool-service/xpu-exporter/server/prometheus.go](/xpu-pool-service/xpu-exporter/server/prometheus.go) | Go | -223 | -23 | -31 | -277 |
| [xpu-pool-service/xpu-exporter/server/prometheus\_test.go](/xpu-pool-service/xpu-exporter/server/prometheus_test.go) | Go | -172 | -5 | -32 | -209 |
| [xpu-pool-service/xpu-exporter/versions/version.go](/xpu-pool-service/xpu-exporter/versions/version.go) | Go | -5 | -6 | -1 | -12 |

[Summary](results.md) / [Details](details.md) / [Diff Summary](diff.md) / Diff Details