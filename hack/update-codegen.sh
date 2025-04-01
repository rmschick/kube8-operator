#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

CODEGEN_SCRIPT="./vendor/k8s.io/code-generator/generate-groups.sh"
GENERATORS="all"
OUTPUT_PACKAGE="kube8-operator/pkg/generated"
APIS_PACKAGE="kube8-operator/pkg/apis"
GROUPS_VERSIONS="collector:v1alpha"
OUTPUT_BASE="."
GO_HEADER_FILE="./hack/boilerplate.go.txt"

chmod +x ${CODEGEN_SCRIPT}
"${CODEGEN_SCRIPT}" "${GENERATORS}" "${OUTPUT_PACKAGE}" "${APIS_PACKAGE}" "${GROUPS_VERSIONS}" \
  --output-base "${OUTPUT_BASE}" \
  --go-header-file "${GO_HEADER_FILE}"

# Copy generated files to the target directory
TARGET_DIR="kube8-operator/pkg"

cp -rf "${OUTPUT_PACKAGE}" "${TARGET_DIR}"

cp -rf "${APIS_PACKAGE}/collector/v1alpha/zz_generated.deepcopy.go" "${TARGET_DIR}/apis/collector/v1alpha"

# Delete the generated code folder
rm -rf "kube8-operator"

