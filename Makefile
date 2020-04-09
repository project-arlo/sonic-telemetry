ifeq ($(GOPATH),)
export GOPATH=/tmp/go
endif
export PATH := $(PATH):$(GOPATH)/bin

INSTALL := /usr/bin/install
DBDIR := /var/run/redis/sonic-db/
GO := /usr/local/go/bin/go
TOP_DIR := $(abspath ..)
GO_MGMT_PATH=$(TOP_DIR)/sonic-mgmt-framework
BUILD_DIR := $(GOPATH)/bin
export CVL_SCHEMA_PATH := $(GO_MGMT_PATH)/src/cvl/schema

SRC_FILES=$(shell find . -name '*.go' | grep -v '_test.go' | grep -v '/tests/')
TEST_FILES=$(wildcard *_test.go)
TELEMETRY_TEST_DIR = $(GO_MGMT_PATH)/build/tests/gnmi_server
TELEMETRY_TEST_BIN = $(TELEMETRY_TEST_DIR)/server.test

all: sonic-telemetry $(TELEMETRY_TEST_BIN)

go.mod:
	/usr/local/go/bin/go mod init github.com/Azure/sonic-telemetry

sonic-telemetry: go.mod
	$(GO) mod vendor
	cp -r $(GO_MGMT_PATH)/vendor/github.com/antchfx ./vendor/github.com/
	cp -r $(GO_MGMT_PATH)/vendor/github.com/openconfig ./vendor/github.com/
	cp -r $(GOPATH)/pkg/mod/github.com/jipanyang/gnxi@v0.0.0-20181221084354-f0a90cca6fd0/* vendor/github.com/jipanyang/gnxi/
	cp -r $(GOPATH)/pkg/mod/golang.org/x/crypto@v0.0.0-20200302210943-78000ba7a073/* vendor/golang.org/x/crypto/
	chmod -R u+w vendor
	patch -d vendor -p0 <patches/gnmi_cli.all.patch
	patch -d vendor -p0 <patches/gnmi_set.patch
	patch -d vendor -p0 <patches/gnmi_get.patch

	$(GO) install -mod=vendor github.com/Azure/sonic-telemetry/telemetry
	$(GO) install -mod=vendor github.com/Azure/sonic-telemetry/dialout/dialout_client_cli
	$(GO) install -mod=vendor github.com/Azure/sonic-telemetry/gnoi_client
	$(GO) install -mod=vendor github.com/jipanyang/gnxi/gnmi_get
	$(GO) install -mod=vendor github.com/jipanyang/gnxi/gnmi_set
	$(GO) install -mod=vendor github.com/openconfig/gnmi/cmd/gnmi_cli

check:
	sudo mkdir -p ${DBDIR}
	sudo cp ./testdata/database_config.json ${DBDIR}
	sudo mkdir -p /usr/models/yang || true
	sudo find $(GO_MGMT_PATH)/models -name '*.yang' -exec cp {} /usr/models/yang/ \;
	-$(GO) test -mod=vendor -v github.com/Azure/sonic-telemetry/gnmi_server
	-$(GO) test -mod=vendor -v github.com/Azure/sonic-telemetry/dialout/dialout_client

clean:
# 	rm -rf vendor
	chmod -f -R u+w $(GOPATH)/pkg || true
	rm -rf $(GOPATH)
	rm -f src

$(TELEMETRY_TEST_BIN): $(TEST_FILES) $(SRC_FILES)
	$(GO) test -mod=vendor -c -cover github.com/Azure/sonic-telemetry/gnmi_server -o $@
	cp -r testdata $(TELEMETRY_TEST_DIR)
	cp -r $(GO_MGMT_PATH)/build/cvl/schema $(TELEMETRY_TEST_DIR)

install:
	$(INSTALL) -D $(BUILD_DIR)/telemetry $(DESTDIR)/usr/sbin/telemetry
	$(INSTALL) -D $(BUILD_DIR)/dialout_client_cli $(DESTDIR)/usr/sbin/dialout_client_cli
	$(INSTALL) -D $(BUILD_DIR)/gnmi_get $(DESTDIR)/usr/sbin/gnmi_get
	$(INSTALL) -D $(BUILD_DIR)/gnmi_set $(DESTDIR)/usr/sbin/gnmi_set
	$(INSTALL) -D $(BUILD_DIR)/gnmi_cli $(DESTDIR)/usr/sbin/gnmi_cli
	$(INSTALL) -D $(BUILD_DIR)/gnoi_client $(DESTDIR)/usr/sbin/gnoi_client

	mkdir -p $(DESTDIR)/usr/bin/
	cp -r $(GO_MGMT_PATH)/build/cvl/schema $(DESTDIR)/usr/sbin
	mkdir -p $(DESTDIR)/usr/models/yang
	find $(GO_MGMT_PATH)/models -name '*.yang' -exec cp {} $(DESTDIR)/usr/models/yang/ \;

deinstall:
	rm $(DESTDIR)/usr/sbin/telemetry
	rm $(DESTDIR)/usr/sbin/dialout_client_cli
	rm $(DESTDIR)/usr/sbin/gnmi_get
	rm $(DESTDIR)/usr/sbin/gnmi_set


