# Copyright 2015-2018 Bret Jordan, All rights reserved.
#
# Use of this source code is governed by an Apache 2.0 license
# that can be found in the LICENSE file in the root of the source tree.

GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_GET=$(GO_CMD) get
GO_INSTALL=$(GO_CMD) install -v
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m


# Binary filename
BINARY=freetaxii
VERSION=0.3.1
BUILD_DIR = srcbuild
BIN_DIR = bin
LOG_DIR = log
DB_DIR = db
ETC_DIR = etc
TEMPLATES_DIR = templates

# The build version that we want to pass in to the application during compile time
BUILD=`git rev-parse HEAD`

# Setup the -ldflags option for go build here, interpolate the variable values 
# LDFLAGS=-ldflags "-X main.Build=$(BUILD)"


# Default target builds FreeTAXII
default:
	@echo "$(OK_COLOR)==> Please run \"make distro\"...$(NO_COLOR)";


# Installs FreeTAXII and copies needed files
distro:
	@echo "$(OK_COLOR)==> Removing Existing Distribution Package...$(NO_COLOR)"; \
	if [ -d $(BUILD_DIR) ] ; then rm -rf $(BUILD_DIR) ; fi

	@echo "$(OK_COLOR)==> Setting Up Distribution Directories...$(NO_COLOR)"; \
	mkdir -p $(BUILD_DIR)/$(BINARY)-$(VERSION)/$(BIN_DIR); \
	mkdir -p $(BUILD_DIR)/$(BINARY)-$(VERSION)/$(DB_DIR); \
	mkdir -p $(BUILD_DIR)/$(BINARY)-$(VERSION)/$(ETC_DIR)/tls; \
	mkdir -p $(BUILD_DIR)/$(BINARY)-$(VERSION)/$(LOG_DIR); \
	mkdir -p $(BUILD_DIR)/$(BINARY)-$(VERSION)/$(TEMPLATES_DIR);

	@echo "$(OK_COLOR)==> Building Application Files...$(NO_COLOR)"; \
	$(GO_BUILD) -v -o $(BUILD_DIR)/$(BINARY)-$(VERSION)/$(BINARY) cmd/freetaxii/freetaxii.go; \
	$(GO_BUILD) -v -o $(BUILD_DIR)/$(BINARY)-$(VERSION)/$(BIN_DIR)/createSqlite3Database cmd/createdb/createSqlite3Database.go; \
	$(GO_BUILD) -v -o $(BUILD_DIR)/$(BINARY)-$(VERSION)/$(BIN_DIR)/verifyconfig cmd/verifyconfig/verifyconfig.go;

	@echo "$(OK_COLOR)==> Copying Needed Files...$(NO_COLOR)"; \
	cp -R cmd/freetaxii/templates/* $(BUILD_DIR)/$(BINARY)-$(VERSION)/$(TEMPLATES_DIR)/; \
	cp cmd/freetaxii/etc/freetaxii.conf $(BUILD_DIR)/$(BINARY)-$(VERSION)/$(ETC_DIR)/; \
	touch $(BUILD_DIR)/$(BINARY)-$(VERSION)/$(LOG_DIR)/$(BINARY).log;

	@echo "$(OK_COLOR)==> Creating Database File...$(NO_COLOR)"; \
	cd $(BUILD_DIR)/$(BINARY)-$(VERSION)/$(BIN_DIR)/; \
	./createSqlite3Database;

	@echo "$(OK_COLOR)==> Creating Tarball...$(NO_COLOR)"; \
	cd $(BUILD_DIR)/; \
	tar -cf $(BINARY)-$(VERSION).tar $(BINARY)-$(VERSION);

	@echo "$(OK_COLOR)==> Compressing Tarball...$(NO_COLOR)"; \
	cd $(BUILD_DIR)/; \
	gzip $(BINARY)-$(VERSION).tar; 


# Clean up the project: delete binaries
clean:
	@echo "$(OK_COLOR)==> Removing Existing Distribution Package $(BINARY)...$(NO_COLOR)"; \
	if [ -d $(BUILD_DIR) ] ; then rm -rf $(BUILD_DIR) ; fi


.PHONY: clean install

