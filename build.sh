#!/usr/bin/env bash

# build folder
buildfolder="build"
if [ ! -d "${buildfolder}" ]; then
  mkdir "${buildfolder}"
else
  rm -fr "${buildfolder}"
fi
nowdate=$(date '+%Y-%m-%d %H:%M:%S')

# logs
sys_log="./${buildfolder}/build.log" #操作日志存放路径

# func of log
function log_warn() {
  echo "[ECHO_BUILD]Warning:${nowdate} $1" >>$sys_log
}

function log_info() {
  echo "[ECHO_BUILD]Info:${nowdate} $1" >>$sys_log
}

function log_err() {
  echo "[ECHO_BUILD]Error:${nowdate} $1" >>$sys_log
}
log_info "build start...."
log_info "build time: ${nowdate}"

# update swagger api
swag init

# set compilation platform
OSARRAY=("darwin" "windows" "linux")
# set build tag
timeTag="-X 'main.BuildTime=$(date '+%Y-%m-%d %H:%M:%S')'"
timeTagCompact="$(date '+%Y%m%d%H%M%S')"
branchName="master"
branchFlag="0000000"
if [ -d ".git" ]; then
  branchName="$(git name-rev --name-only HEAD)"
  commitHash="$(git rev-parse --short HEAD)"
else
  log_warn "not a git repository!"
fi
log_info "build branch: ${branchName}"
log_info "commit hash: ${branchFlag}"

# set golang vars
branchFlag="-X main.GitBranch=${branchName}"
commitFlag="-X main.CommitId=${commitHash}"
goVersionFlag="-X main.GoVersion=$(go version | awk '{print $3}')"
ldflags="-s -w ${timeTag} ${branchFlag} ${commitFlag} ${goVersionFlag} "
log_info "go version: $(go version)"

# build
appName=$(basename "$(pwd)")
for os in ${OSARRAY[*]}; do
  buildname="./${buildfolder}/${appName}_${os}_amd64_${branchName}_${timeTagCompact}"
  if [ "${os}" == "windows" ]; then
    buildname="${buildname}.exe"
  fi
  log_info "building-> ${buildname}"
  CGO_ENABLED=0 GOOS="${os}" GOARCH=amd64 go build -ldflags "${ldflags}" -o "${buildname}"
  # check
  if [ ! -f "${buildname}" ]; then
    log_err "build failed-> ${buildname}"
  fi
done
cp config_default.toml "./${buildfolder}/config.toml"
log_info "building complete..."

# clean
log_info "clean..."
tar -zcf "${appName}_amd64_${timeTagCompact}.tar.gz" "./${buildfolder}"
rm -fr "./${buildfolder}"
