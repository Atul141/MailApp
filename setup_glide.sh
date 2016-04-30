#!/bin/bash

set -o pipefail -o errexit

CDTOOLS_RUN_PATH="${GOPATH}/bin"
GLIDE_VERSION="0.8.3"

GITHUB_SRC_PATH="${GOPATH}/src/github.com"
GLIDE_INSTALL_SRC_PATH="${GITHUB_SRC_PATH}/Masterminds"
GLIDE_SRC_PATH="${GLIDE_INSTALL_SRC_PATH}/glide"
GLIDE_GITHUB_REPO="https://github.com/Masterminds/glide.git"

if [[ -n "${GO_PIPELINE_NAME}" ]]; then
	CDTOOLS_RUN_PATH="/var/go/cdtools"
fi

mkdir -p "${GITHUB_SRC_PATH}"
cd "${GITHUB_SRC_PATH}"

if [[ ! -d "${GLIDE_INSTALL_SRC_PATH}" ]]; then 
	mkdir -p "Masterminds"
fi

cd "${GLIDE_INSTALL_SRC_PATH}"

if [[ -d "${GLIDE_SRC_PATH}" ]]; then
  rm -rf "${GLIDE_SRC_PATH}"
fi

git clone "${GLIDE_GITHUB_REPO}"

cd "${GLIDE_SRC_PATH}"
git checkout tags/"${GLIDE_VERSION}" -b "${GLIDE_VERSION}"

echo "Building Glide from source"

make bootstrap
make build

if [[ "$?" != 0 ]]; then
	echo "Failed building glide"
	exit 1
fi

if [[ -f "${GLIDE_SRC_PATH}/glide" ]]; then
  echo "Copying glide executable to ${CDTOOLS_RUN_PATH}"
  cp glide "${CDTOOLS_RUN_PATH}"
else
  echo "Glide executable not generated"
  exit 1
fi

echo "Installation of Glide is complete"
