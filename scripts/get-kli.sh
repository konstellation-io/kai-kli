#!/usr/bin/env bash

# Copyright (2023) konstellation-io.
#
# Licensed under the MIT License.
# See https://github.com/konstellation-io/kai-kli/blob/develop/LICENSE for details

# The install script is based off of the Apache-licensed script from Helm,
# the package manager for Kubernetes: https://github.com/helm/helm/blob/main/scripts/get-helm-3

: ${BINARY_NAME:="kli"}
: ${USE_SUDO:="true"}
: ${DEBUG:="false"}
: ${VERIFY_CHECKSUM:="true"}
: ${KLI_INSTALL_DIR:="/usr/local/bin"}

HAS_CURL="$(type "curl" &> /dev/null && echo true || echo false)"
HAS_WGET="$(type "wget" &> /dev/null && echo true || echo false)"
HAS_OPENSSL="$(type "openssl" &> /dev/null && echo true || echo false)"
HAS_GIT="$(type "git" &> /dev/null && echo true || echo false)"

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
    aarch64) ARCH="arm64";;
    x86_64) ARCH="amd64";;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Minimalist GNU for Windows
    mingw*|cygwin*) OS='windows';;
  esac
}

# runs the given command as root (detects if we are root already)
runAsRoot() {
  if [ $EUID -ne 0 -a "$USE_SUDO" = "true" ]; then
    sudo "${@}"
  else
    "${@}"
  fi
}

# verifySupported checks that the os/arch combination is supported for
# binary builds, as well whether or not necessary tools are present.
verifySupported() {
  local supported="darwin-amd64\ndarwin-arm64\nlinux-amd64\nlinux-arm64\nwindows-amd64"
  if ! echo "${supported}" | grep -q "${OS}-${ARCH}"; then
    echo "No prebuilt binary for ${OS}-${ARCH}."
    echo "To build from source, go to https://github.com/konstellation-io/kai-kli"
    exit 1
  fi

  if [ "${HAS_CURL}" != "true" ] && [ "${HAS_WGET}" != "true" ]; then
    echo "Either curl or wget is required"
    exit 1
  fi

  if [ "${VERIFY_CHECKSUM}" == "true" ] && [ "${HAS_OPENSSL}" != "true" ]; then
    echo "In order to verify checksum, openssl must first be installed."
    echo "Please install openssl or set VERIFY_CHECKSUM=false in your environment."
    exit 1
  fi

  if [ "${HAS_GIT}" != "true" ]; then
    echo "[WARNING] Could not find git. It is required for plugin installation."
  fi
}

# checkDesiredVersion checks if the desired version is available.
checkDesiredVersion() {
  if [ "x$DESIRED_VERSION" == "x" ]; then
    # Get tag from release URL
    local latest_release_url="https://api.github.com/repos/konstellation-io/kai-kli/releases/latest"
    local latest_release_response=""
    if [ "${HAS_CURL}" == "true" ]; then
      latest_release_response=$( curl -L --silent --show-error --fail "$latest_release_url" 2>&1 || true )
    elif [ "${HAS_WGET}" == "true" ]; then
      latest_release_response=$( wget "$latest_release_url" -O - 2>&1 || true )
    fi
    TAG=$( echo "$latest_release_response" | grep '"tag_name"' | sed -E 's/.*"(v[0-9\.]+)".*/\1/g' )
    if [ "x$TAG" == "x" ]; then
      printf "Could not retrieve the latest release tag information from %s: %s\n" "${latest_release_url}" "${latest_release_response}"
      exit 1
    fi
  else
    TAG=$DESIRED_VERSION
  fi
}

# checkKliInstalledVersion checks which version of kli is installed and
# if it needs to be changed.
checkKliInstalledVersion() {
  if [[ -f "${KLI_INSTALL_DIR}/${BINARY_NAME}" ]]; then
    local version=$("${KLI_INSTALL_DIR}/${BINARY_NAME}" version | awk '{print $3}')
    if [[ "$version" == "${TAG#v}" ]]; then
      echo "Kli ${version} is already ${DESIRED_VERSION:-latest}"
      return 0
    else
      echo "Kli ${TAG} is available. Changing from version ${version}."
      return 1
    fi
  else
    return 1
  fi
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  KLI_DIST="kai-kli_${TAG#v}_${OS}-${ARCH}.tar.gz"
  KLI_SUM="kai-kli_${TAG#v}_checksums.txt"
  DOWNLOAD_URL="https://github.com/konstellation-io/kai-kli/releases/download/$TAG/$KLI_DIST"
  CHECKSUM_URL="https://github.com/konstellation-io/kai-kli/releases/download/$TAG/$KLI_SUM"
  KLI_TMP_ROOT="$(mktemp -dt kli-installer-XXXXXX)"
  KLI_TMP_FILE="$KLI_TMP_ROOT/$KLI_DIST"
  KLI_SUM_FILE="$KLI_TMP_ROOT/$KLI_SUM"
  echo "Downloading $DOWNLOAD_URL"
  if [ "${HAS_CURL}" == "true" ]; then
    curl -SsL "$CHECKSUM_URL" -o "$KLI_SUM_FILE"
    curl -SsL "$DOWNLOAD_URL" -o "$KLI_TMP_FILE"
  elif [ "${HAS_WGET}" == "true" ]; then
    wget -q -O "$KLI_SUM_FILE" "$CHECKSUM_URL"
    wget -q -O "$KLI_TMP_FILE" "$DOWNLOAD_URL"
  fi
}

# verifyFile verifies the SHA256 checksum of the binary package
# (depending on settings in environment).
verifyFile() {
  if [ "${VERIFY_CHECKSUM}" == "true" ]; then
    verifyChecksum
  fi
}

# installFile installs the Kli binary.
installFile() {
  KLI_TEMP="$KLI_TMP_ROOT/$BINARY_NAME"
  mkdir -p "$KLI_TEMP"
  tar xf "$KLI_TMP_FILE" -C "$KLI_TEMP"
  KLI_TEMP_BIN="$KLI_TEMP/kai-kli_${TAG#v}_$OS-$ARCH/bin/kli"
  echo "Preparing to install $BINARY_NAME into ${KLI_INSTALL_DIR}"
  runAsRoot cp "$KLI_TEMP_BIN" "$KLI_INSTALL_DIR/$BINARY_NAME"
  echo "$BINARY_NAME installed into $KLI_INSTALL_DIR/$BINARY_NAME"
}

# verifyChecksum verifies the SHA256 checksum of the binary package.
verifyChecksum() {
  printf "Verifying checksum... "
  local sum=$(openssl sha1 -sha256 ${KLI_TMP_FILE} | awk '{print $2}')
  local expected_sum=$(cat ${KLI_SUM_FILE} | grep "${OS}-${ARCH}.tar.gz" | awk '{print $1}')
  if [ "$sum" != "$expected_sum" ]; then
    echo "SHA sum of ${KLI_TMP_FILE} does not match. Aborting."
    exit 1
  fi
  echo "Done."
}

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    if [[ -n "$INPUT_ARGUMENTS" ]]; then
      echo "Failed to install $BINARY_NAME with the arguments provided: $INPUT_ARGUMENTS"
      help
    else
      echo "Failed to install $BINARY_NAME"
    fi
    echo -e "\tFor support, go to https://github.com/konstellation-io/kai-kli."
  fi
  cleanup
  exit $result
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
  set +e
  KAI="$(command -v $BINARY_NAME)"
  if [ "$?" = "1" ]; then
    echo "$BINARY_NAME not found. Is $KLI_INSTALL_DIR on your "'$PATH?'
    exit 1
  fi
  set -e
}

# help provides possible cli installation arguments
help () {
  echo "Accepted cli arguments are:"
  echo -e "\t[--help|-h ] ->> prints this help"
  echo -e "\t[--version|-v <desired_version>] . When not defined it fetches the latest release from GitHub"
  echo -e "\te.g. --version v2.0.0"
  echo -e "\t[--no-sudo]  ->> install without sudo"
}

# cleanup temporary files
cleanup() {
  if [[ -d "${KLI_TMP_ROOT:-}" ]]; then
    rm -rf "$KLI_TMP_ROOT"
  fi
}

# Execution

#Stop execution on any error
trap "fail_trap" EXIT
set -e

# Set debug if desired
if [ "${DEBUG}" == "true" ]; then
  set -x
fi

# Parsing input arguments (if any)
export INPUT_ARGUMENTS="${@}"
set -u
while [[ $# -gt 0 ]]; do
  case $1 in
    '--version'|-v)
       shift
       if [[ $# -ne 0 ]]; then
           export DESIRED_VERSION="${1}"
           if [[ "$1" != "v"* ]]; then
               echo "Expected version arg ('${DESIRED_VERSION}') to begin with 'v', fixing..."
               export DESIRED_VERSION="v${1}"
           fi
       else
           echo -e "Please provide the desired version. e.g. --version v3.0.0 or -v canary"
           exit 0
       fi
       ;;
    '--no-sudo')
       USE_SUDO="false"
       ;;
    '--help'|-h)
       help
       exit 0
       ;;
    *) exit 1
       ;;
  esac
  shift
done
set +u

initArch
initOS
verifySupported
checkDesiredVersion
if ! checkKliInstalledVersion; then
  downloadFile
  verifyFile
  installFile
fi
testVersion
cleanup
