#!/bin/sh

set -euo pipefail

# increase script verbosity with -x
# set -euxo pipefail

# version to fetch
VERSION="v1.0.1"

# download binary and plugin configuration
curl -L -o cleanup_${VERSION}.tgz https://github.com/ashleyschuett/kubeconfig-cleanup/releases/download/${VERSION}/$BINARY
curl -LO https://raw.githubusercontent.com/ashleyschuett/kubeconfig-cleanup/${VERSION}/plugin.yaml

# extract binary
tar zxvf cleanup_${VERSION}.tgz cleanup
rm cleanup_${VERSION}.tgz

# create corresponding plugin folder and move binary/configuration
mkdir -p ~/.kube/plugins/kubeconfig-cleanup
mv cleanup plugin.yaml ~/.kube/plugins/kubeconfig-cleanup
