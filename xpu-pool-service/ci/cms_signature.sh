#!/bin/bash
# Copyright Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
# 构建签名脚本
set -e

current_dir=$(
    cd "$(dirname "$0")" || exit 1
    pwd
)
workspace=$(dirname "${current_dir}")
pkg_path=$1
signature_jar=$(find /opt/buildtools/ -name signature.jar)

if [ ! -d "${workspace}/CI" ]; then
    mkdir -p "${workspace}/CI"
fi

function gen_list() {
    for file in "$1"/*; do
        if [ -d "$file" ]; then
            gen_list "$file"
        else
            echo "$file" is file
            if [ "$(basename "$file")"x != listx ]; then
                cat <<EOF >> "${pkg_path}/list"
Name: ${file##*/pkg_path}
SHA256-Digest: $(sha256sum "$file" | awk '{print $1}')
EOF
            fi
        done
    }
}

function gen_signature_xml() {
    cat << EOF > "${workspace}/CI/signconf_cms.xml"
<?xml version="1.0" encoding="UTF-8"?>
<!-- 由产品CI配置此文件，供私有构建、团队构建、发布构建等各级工程共享 -->
<signtasks>
  <signtask name="linux_single">
    <alias>CMS_Computing_RSA2048_CN_20220810_Huawei</alias>
    <fileset path="${pkg_path}">
      <includes>**/list</includes>
    </fileset>
    <crfiles>${pkg_path}/list.cms.cnk/crfiles</crfiles>
    <hashtype>2</hashtype>
    <proxylist>10.29.154.209:12056</proxylist>
    <signaturestandard>5</signaturestandard>
    <productlineid>849944</productlineid>
    <versionid>260181123</versionid>
    <padmode>1</padmode>
  </signtask>
</signtasks>
EOF
}

cd "${pkg_path}"
cat <<EOF >"${pkg_path}/list"
Manifest Version: 1.0
Create By: Huawei Technology Inc.
EOF

gen_list "${pkg_path}"
gen_signature_xml
java -jar "${signature_jar}" "${workspace}/CI/signconf_cms.xml"