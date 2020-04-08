BIN            = gitcomm
BUILD         ?= $(shell git rev-parse --short HEAD)
BUILD_DATE    ?= $(shell git log -1 --format=%ai)
BUILD_BRANCH  ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILD_VERSION ?= $(shell git describe --always --tags)
BUILD_TAGS    ?=
GOPATH        ?= $(shell go env GOPATH)

BASEPATH = github.com/karantin2020/gitcomm/version

export GO111MODULE := off

# Build-time Go variables
appVersion     = ${BASEPATH}.version
gitBranch      = ${BASEPATH}.gitBranch
lastCommitSHA  = ${BASEPATH}.lastCommitSHA
lastCommitTime = ${BASEPATH}.lastCommitTime

BUILD_FLAGS   ?= -ldflags '-s -w -X ${lastCommitSHA}=${BUILD} -X "${lastCommitTime}=${BUILD_DATE}" -X "${appVersion}=${BUILD_VERSION}" -X ${gitBranch}=${BUILD_BRANCH}'

build:
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -v -o cmd/${BIN} ./cmd