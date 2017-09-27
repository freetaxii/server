# Copyright 2017 Bret Jordan, All rights reserved.
#
# Use of this source code is governed by an Apache 2.0 license
# that can be found in the LICENSE file in the root of the source
# tree.

GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_INSTALL=$(GO_CMD) install -v
GO_CLEAN=$(GO_CMD) clean
GO_LINT=golint
GO_VET=$(GO_CMD) vet
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m


# Binary filename
BINARY=freetaxii


# Variables we want to pass in to the application during compile time
VERSION=0.0.1
BUILD=`git rev-parse HEAD`

# Setup the -ldflags option for go build here, interpolate the variable values 
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"


# Default target builds FreeTAXII
default:
	@echo "${OK_COLOR}==> Building ${BINARY}...${NO_COLOR}"; \
	${GO_BUILD} ${LDFLAGS} -o ${BINARY}

# Build a version specifically for Darwin 64-bit
darwin:
	@echo "${OK_COLOR}==> Building ${BINARY} for Darwin...${NO_COLOR}"; \
	GOOS=darwin GOARCH=amd64 ${GO_BUILD} ${LDFLAGS} -o ${BINARY}-darwin-amd64

# Build a version specificatlly for Linux 64-bit
linux64:
	@echo "${OK_COLOR}==> Building ${BINARY} for Linux64...${NO_COLOR}"; \
	GOOS=linux GOARCH=amd64 ${GO_BUILD} ${LDFLAGS} -o ${BINARY}-linux-amd64	

# Installs FreeTAXII and copies needed files
install:
	@echo "${OK_COLOR}==> Installing ${BINARY}...${NO_COLOR}"; \
	${GO_INSTALL} ${LDFLAGS}

# Clean up the project: delete binaries
clean:
	@echo "${OK_COLOR}==> Cleaning ${BINARY}...${NO_COLOR}"; \
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

# Run the go vet on the go code
vet:
	@echo "${OK_COLOR}==> Vetting ${BINARY}...${NO_COLOR}"; \
	${GO_VET} ./...

# Run the go linter on the go code
lint:
	@echo "${OK_COLOR}==> Linting ${BINARY}...${NO_COLOR}"; \
	${GO_LINT} ./...

.PHONY: clean install vet lint

