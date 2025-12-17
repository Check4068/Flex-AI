module huawei.com/vxpu-device-plugin

go 1.22.1

replace (
	github.com/fsnotify/fsnotify => github.com/fsnotify/fsnotify v1.7.0
	golang.org/x/net => golang.org/x/net v0.20.0
	golang.org/x/sync => golang.org/x/sync v0.6.0
	golang.org/x/sys => golang.org/x/sys v0.17.0
	golang.org/x/text => golang.org/x/text v0.14.0
)

require (
	github.com/agiledragon/gomonkey/v2 v2.8.0
	github.com/fsnotify/fsnotify v1.6.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.9.0
	golang.org/x/net v0.20.0
	golang.org/x/sync v0.6.0
	golang.org/x/sys v0.17.0
	golang.org/x/text v0.14.0
	gopkg.in/yaml.v2 v2.4.0
	huawei.com/npu-exporter/v6 v6.0.0-RC2.b001
	k8s.io/api v0.31.1
	k8s.io/apimachinery v0.31.1
	k8s.io/client-go v0.31.1
	k8s.io/kubelet v0.31.1
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20210407192527-94a9f03dee38 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0-20181201164244-5d4384ee4fb2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/tools v0.18.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240701130421-f6361c86f094 // indirect
	google.golang.org/grpc v1.57.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.130.0 // indirect
	k8s.io/kube-openapi v0.0.0-20240228011516-7ddd3763d340 // indirect
	k8s.io/utils v0.0.0-202402280113017-18e5e09b52bc // indirect
	sigs.k8s.io/json v0.0.0-20221106044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)