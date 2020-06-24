# Copyright (c) 2020 ATS, Inc. All rights reserved.
#
# This Makefile requires the following dependencies on the build system:
# - go
#
SB_TOP := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
PARENT_DIR := $(abspath $(SB_TOP)/../)/
SHELL=/bin/bash -o pipefail
OUT_DIR := $(abspath $(SB_TOP)/build/_output/)
OPERATOR_BIN := $(abspath $(SB_TOP)/build/_output/bin/tf-operator)

OPERATOR_VERSION := v0.0.1

CGO_OPT := "CGO_ENABLED=0"
GIT_REPO := "github.com/atsgen/tf-operator/"
MANAGER_DIR := "$(GIT_REPO)cmd/manager"

GC_FLAGS := "-gcflags=all=-trimpath=$(SB_TOP)"
ASM_FLAGS := "-asmflags=all=-trimpath=$(SB_TOP)"

all: build image

.PHONY: build
build:
	@export $(CGO_OPT); go build -o $(OPERATOR_BIN) $(GC_FLAGS) $(ASM_FLAGS) -mod=vendor $(MANAGER_DIR)
	@echo "Build done for tf-operator: $(OPERATOR_BIN)"

.PHONY: image
image:
	docker build --build-arg OPERATOR_VERSION=$(OPERATOR_VERSION) -t atsgen/tf-operator:$(OPERATOR_VERSION) -f build/Dockerfile .
	docker build --build-arg OPERATOR_VERSION=$(OPERATOR_VERSION) -t atsgen/tf-operator:$(OPERATOR_VERSION)-ubi -f build/Dockerfile.ubi .

.PHONY: clean
clean:
	@rm -rf $(OUT_DIR)
#go build -o /root/contrail/tf-operator/build/_output/bin/tf-operator -gcflags all=-trimpath=/root/contrail -asmflags all=-trimpath=/root/contrail -mod=vendor github.com/atsgen/tf-operator/cmd/manager

