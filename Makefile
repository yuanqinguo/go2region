export VERSION=1.0.0
export ENV=prod
export PROJECT=goip2region

TOPDIR=$(shell pwd)
SOURCE_BINARY_DIR=$(TOPDIR)/bin
SOURCE_BINARY_FILE=$(SOURCE_BINARY_DIR)/$(PROJECT)
SOURCE_MAIN_FILE=main.go

BUILD_TIME=`date +%Y%m%d%H%M%S`
BUILD_FLAG=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

OBJTAR=$(PROJECT).tar.gz

all: build pack
	@echo "ALL DONE"
	@echo "Program:       "  $(PROJECT)
	@echo "Version:       "  $(VERSION)
	@echo "Env:          "  $(ENV)

build:
	@echo "start go build...." $(TOPDIR)
	@rm -rf $(SOURCE_BINARY_DIR)/*
	@go build $(BUILD_FLAG) -o $(SOURCE_BINARY_FILE) $(SOURCE_MAIN_FILE)

pack:
	@echo "\n\rpacking...."
	@tar czvf $(OBJTAR) -C $(SOURCE_BINARY_DIR) .
