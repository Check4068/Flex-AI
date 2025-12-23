# Diff Details

Date : 2025-12-08 15:35:42

Directory c:\\Users\\admin\\IdeaProjects\\untitled\\scheduler\\xpu-pool-service

Total : 43 files,  2923 codes, 277 comments, 337 blanks, all 3537 lines

[Summary](results.md) / [Details](details.md) / [Diff Summary](diff.md) / Diff Details

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [xpu-pool-service/CMakeLists.txt](/xpu-pool-service/CMakeLists.txt) | CMake | 19 | 0 | 7 | 26 |
| [xpu-pool-service/GPU-device-plugin/Makefile](/xpu-pool-service/GPU-device-plugin/Makefile) | Makefile | 34 | 2 | 6 | 42 |
| [xpu-pool-service/GPU-device-plugin/client/client.go](/xpu-pool-service/GPU-device-plugin/client/client.go) | Go | 53 | 13 | 9 | 75 |
| [xpu-pool-service/GPU-device-plugin/client/client\_test.go](/xpu-pool-service/GPU-device-plugin/client/client_test.go) | Go | 146 | 39 | 17 | 202 |
| [xpu-pool-service/GPU-device-plugin/cmd/main.go](/xpu-pool-service/GPU-device-plugin/cmd/main.go) | Go | 97 | 38 | 25 | 160 |
| [xpu-pool-service/GPU-device-plugin/cmd/main\_test.go](/xpu-pool-service/GPU-device-plugin/cmd/main_test.go) | Go | 347 | 2 | 9 | 358 |
| [xpu-pool-service/GPU-device-plugin/examples/vgpu-exceed.yml](/xpu-pool-service/GPU-device-plugin/examples/vgpu-exceed.yml) | YAML | 14 | 1 | 0 | 15 |
| [xpu-pool-service/GPU-device-plugin/examples/vgpu-test.yml](/xpu-pool-service/GPU-device-plugin/examples/vgpu-test.yml) | YAML | 25 | 0 | 0 | 25 |
| [xpu-pool-service/GPU-device-plugin/examples/vgpu-whole-card.yml](/xpu-pool-service/GPU-device-plugin/examples/vgpu-whole-card.yml) | YAML | 13 | 0 | 0 | 13 |
| [xpu-pool-service/GPU-device-plugin/go.mod](/xpu-pool-service/GPU-device-plugin/go.mod) | Go Module File | 63 | 0 | 4 | 67 |
| [xpu-pool-service/GPU-device-plugin/go.sum](/xpu-pool-service/GPU-device-plugin/go.sum) | Go Checksum File | 177 | 0 | 0 | 177 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/api.pb.go](/xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/api.pb.go) | Go | 289 | 18 | 42 | 349 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/api\_grpc.pb.go](/xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/api_grpc.pb.go) | Go | 108 | 22 | 18 | 148 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/service\_impl.go](/xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/service_impl.go) | Go | 422 | 21 | 33 | 476 |
| [xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/service\_impl\_test.go](/xpu-pool-service/GPU-device-plugin/pkg/api/runtime/service/service_impl_test.go) | Go | 233 | 4 | 22 | 259 |
| [xpu-pool-service/GPU-device-plugin/watchers/watchers.go](/xpu-pool-service/GPU-device-plugin/watchers/watchers.go) | Go | 24 | 6 | 7 | 37 |
| [xpu-pool-service/GPU-device-plugin/watchers/watchers\_test.go](/xpu-pool-service/GPU-device-plugin/watchers/watchers_test.go) | Go | 54 | 4 | 10 | 68 |
| [xpu-pool-service/README.md](/xpu-pool-service/README.md) | Markdown | 0 | 0 | 1 | 1 |
| [xpu-pool-service/ci/VersionSet.xml](/xpu-pool-service/ci/VersionSet.xml) | XML | 1 | 0 | 0 | 1 |
| [xpu-pool-service/ci/app\_define.json](/xpu-pool-service/ci/app_define.json) | JSON | 13 | 0 | 0 | 13 |
| [xpu-pool-service/ci/at/at\_deploy.sh](/xpu-pool-service/ci/at/at_deploy.sh) | Shell Script | 23 | 3 | 5 | 31 |
| [xpu-pool-service/ci/at/at\_deploy.yml](/xpu-pool-service/ci/at/at_deploy.yml) | YAML | 48 | 0 | 5 | 53 |
| [xpu-pool-service/ci/build.sh](/xpu-pool-service/ci/build.sh) | Shell Script | 80 | 4 | 18 | 102 |
| [xpu-pool-service/ci/build.yml](/xpu-pool-service/ci/build.yml) | YAML | 44 | 0 | 4 | 48 |
| [xpu-pool-service/ci/buildinfo.sh](/xpu-pool-service/ci/buildinfo.sh) | Shell Script | 11 | 3 | 2 | 16 |
| [xpu-pool-service/ci/cmc/openSource\_x86.xml](/xpu-pool-service/ci/cmc/openSource_x86.xml) | XML | 24 | 0 | 0 | 24 |
| [xpu-pool-service/ci/cmc/upload\_cmc.xml](/xpu-pool-service/ci/cmc/upload_cmc.xml) | XML | 24 | 0 | 0 | 24 |
| [xpu-pool-service/ci/cms\_signature.sh](/xpu-pool-service/ci/cms_signature.sh) | Shell Script | 55 | 3 | 6 | 64 |
| [xpu-pool-service/ci/dependency.xml](/xpu-pool-service/ci/dependency.xml) | XML | 7 | 1 | 0 | 8 |
| [xpu-pool-service/ci/hwp7s\_signature.sh](/xpu-pool-service/ci/hwp7s_signature.sh) | Shell Script | 37 | 4 | 3 | 44 |
| [xpu-pool-service/ci/opensource.xml](/xpu-pool-service/ci/opensource.xml) | XML | 6 | 0 | 0 | 6 |
| [xpu-pool-service/ci/xpu\_pool/build\_x86.yml](/xpu-pool-service/ci/xpu_pool/build_x86.yml) | YAML | 65 | 0 | 5 | 70 |
| [xpu-pool-service/ci/xpu\_pool/build\_xpu\_package.sh](/xpu-pool-service/ci/xpu_pool/build_xpu_package.sh) | Shell Script | 107 | 3 | 13 | 123 |
| [xpu-pool-service/ci/xpu\_pool/xpu\_docker\_build/acl\_client/Dockerfile](/xpu-pool-service/ci/xpu_pool/xpu_docker_build/acl_client/Dockerfile) | Docker | 5 | 2 | 3 | 10 |
| [xpu-pool-service/ci/xpu\_pool/xpu\_docker\_build/cuda\_client/Dockerfile](/xpu-pool-service/ci/xpu_pool/xpu_docker_build/cuda_client/Dockerfile) | Docker | 5 | 2 | 3 | 10 |
| [xpu-pool-service/ci/xpu\_pool/xpu\_docker\_build/exporter/Dockerfile](/xpu-pool-service/ci/xpu_pool/xpu_docker_build/exporter/Dockerfile) | Docker | 20 | 4 | 6 | 30 |
| [xpu-pool-service/ci/xpu\_pool/xpu\_docker\_build/gpu\_device\_plugin/Dockerfile](/xpu-pool-service/ci/xpu_pool/xpu_docker_build/gpu_device_plugin/Dockerfile) | Docker | 10 | 4 | 5 | 19 |
| [xpu-pool-service/ci/xpu\_pool/xpu\_docker\_build/npu\_device\_plugin/Dockerfile](/xpu-pool-service/ci/xpu_pool/xpu_docker_build/npu_device_plugin/Dockerfile) | Docker | 5 | 10 | 4 | 19 |
| [xpu-pool-service/client\_update/acl-client-update.sh](/xpu-pool-service/client_update/acl-client-update.sh) | Shell Script | 60 | 28 | 15 | 103 |
| [xpu-pool-service/client\_update/cuda-client-update.sh](/xpu-pool-service/client_update/cuda-client-update.sh) | Shell Script | 68 | 35 | 12 | 115 |
| [xpu-pool-service/direct/CMakeLists.txt](/xpu-pool-service/direct/CMakeLists.txt) | CMake | 21 | 0 | 4 | 25 |
| [xpu-pool-service/direct/acl/CMakeLists.txt](/xpu-pool-service/direct/acl/CMakeLists.txt) | CMake | 33 | 0 | 7 | 40 |
| [xpu-pool-service/direct/make\_lib\_original.sh](/xpu-pool-service/direct/make_lib_original.sh) | Shell Script | 33 | 1 | 7 | 41 |

[Summary](results.md) / [Details](details.md) / [Diff Summary](diff.md) / Diff Details