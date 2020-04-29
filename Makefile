ifeq ($(GOPATH),)
	export GOPATH=/tmp/go
endif
export PATH := $(PATH):$(GOPATH)/bin

INSTALL := /usr/bin/install
DBDIR := /var/run/redis/sonic-db/
ifeq ($(GO),)
	GO := /usr/local/go/bin/go 
	export GO
endif

TOP_DIR := $(abspath ..)
MGMT_COMMON_DIR := $(TOP_DIR)/sonic-mgmt-common
BUILD_DIR := build/bin
export CVL_SCHEMA_PATH := $(MGMT_COMMON_DIR)/cvl/schema
export GOBIN := $(abspath $(BUILD_DIR))

SRC_FILES=$(shell find . -name '*.go' | grep -v '_test.go' | grep -v '/tests/')
TEST_FILES=$(wildcard *_test.go)
TELEMETRY_TEST_DIR = build/tests/gnmi_server
TELEMETRY_TEST_BIN = $(TELEMETRY_TEST_DIR)/server.test

all: sonic-telemetry $(TELEMETRY_TEST_BIN)

go.mod:
	$(GO) mod init github.com/Azure/sonic-telemetry

<<<<<<< HEAD
deps: $(BUILD_DIR)/.deps

$(BUILD_DIR)/.deps: $(MAKEFILE_LIST)
	GOPATH=$(GO_DEP_PATH) $(GO) get -u  github.com/Workiva/go-datastructures/queue
	GOPATH=$(GO_DEP_PATH) $(GO) get -u github.com/openconfig/goyang
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/openconfig/ygot/ygot
	GOPATH=$(GO_DEP_PATH) $(GO) get -u github.com/golang/glog
	GOPATH=$(GO_DEP_PATH) $(GO) get -u github.com/go-redis/redis
	GOPATH=$(GO_DEP_PATH) $(GO) get -u  github.com/c9s/goprocinfo/linux
	GOPATH=$(GO_DEP_PATH) $(GO) get -u  github.com/golang/protobuf/proto
	GOPATH=$(GO_DEP_PATH) $(GO) get -d  github.com/openconfig/gnmi/proto/gnmi
	GOPATH=$(GO_DEP_PATH) $(GO) get -u  golang.org/x/net/context
	GOPATH=$(GO_DEP_PATH) $(GO) get -u  google.golang.org/grpc
	GOPATH=$(GO_DEP_PATH) $(GO) get -u google.golang.org/grpc/credentials
	GOPATH=$(GO_DEP_PATH) $(GO) get -u github.com/google/gnxi/utils
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/jipanyang/gnxi/utils/xpath
	cd $(GO_DEP_PATH)/src/github.com/openconfig/gnmi/proto/gnmi; git reset --hard HEAD;git clean -f -d;git checkout e7106f7f5493a9fa152d28ab314f2cc734244ed8 2>/dev/null ; true; \
  GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/openconfig/gnmi/proto/gnmi
	GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/jipanyang/gnxi/utils/xpath
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/jipanyang/gnmi/client/gnmi
	GOPATH=$(GO_DEP_PATH) $(GO) get -u github.com/xeipuuv/gojsonschema
	GOPATH=$(GO_DEP_PATH) $(GO) get -u github.com/openconfig/gnoi/system
	GOPATH=$(GO_DEP_PATH) $(GO) get -u github.com/msteinert/pam
	GOPATH=$(GO_DEP_PATH) $(GO) get -u github.com/dgrijalva/jwt-go
	GOPATH=$(GO_DEP_PATH) $(GO) get -u gopkg.in/godbus/dbus.v5
	GOPATH=$(GO_DEP_PATH) $(GO) get -u github.com/gogo/protobuf/gogoproto
	
	cd $(GO_DEP_PATH)/src/github.com/openconfig/ygot/; git reset --hard HEAD;git clean -f -d;git checkout 724a6b18a9224343ef04fe49199dfb6020ce132a 2>/dev/null ; true; \
GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/openconfig/ygot/ygot
	cd $(GO_DEP_PATH)/src/github.com/jipanyang/gnmi/client/gnmi; git reset --hard HEAD;git clean -f -d;git checkout cb4d464fa018c29eadab98281448000bee4dcc3d 2>/dev/null ; true; \
GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/jipanyang/gnmi/client/gnmi
=======
sonic-telemetry: $(MAKEFILE_LIST) go.mod
	$(GO) mod vendor
	$(MGMT_COMMON_DIR)/patches/apply.sh vendor
	cp -r $(GOPATH)/pkg/mod/github.com/jipanyang/gnxi@v0.0.0-20181221084354-f0a90cca6fd0/* vendor/github.com/jipanyang/gnxi/
	cp -r $(GOPATH)/pkg/mod/golang.org/x/crypto@v0.0.0-20200302210943-78000ba7a073/* vendor/golang.org/x/crypto/
	chmod -R u+w vendor
	patch -d vendor -p0 < patches/gnmi_cli.all.patch
	patch -d vendor -p0 < patches/gnmi_set.patch
	patch -d vendor -p0 < patches/gnmi_get.patch
>>>>>>> broadcom_sonic_3.x_share
	
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
	sudo find $(MGMT_COMMON_DIR)/models -name '*.yang' -exec cp {} /usr/models/yang/ \;
	-$(GO) test -mod=vendor -v github.com/Azure/sonic-telemetry/gnmi_server
	-$(GO) test -mod=vendor -v github.com/Azure/sonic-telemetry/dialout/dialout_client

clean:
	rm -rf build
	rm -rf vendor

$(TELEMETRY_TEST_BIN): $(TEST_FILES) $(SRC_FILES)
	mkdir -p $(@D)
	cp -r testdata $(@D)/
	$(GO) test -mod=vendor -c -cover github.com/Azure/sonic-telemetry/gnmi_server -o $@

install:
	$(INSTALL) -D $(BUILD_DIR)/telemetry $(DESTDIR)/usr/sbin/telemetry
	$(INSTALL) -D $(BUILD_DIR)/dialout_client_cli $(DESTDIR)/usr/sbin/dialout_client_cli
	$(INSTALL) -D $(BUILD_DIR)/gnmi_get $(DESTDIR)/usr/sbin/gnmi_get
	$(INSTALL) -D $(BUILD_DIR)/gnmi_set $(DESTDIR)/usr/sbin/gnmi_set
	$(INSTALL) -D $(BUILD_DIR)/gnmi_cli $(DESTDIR)/usr/sbin/gnmi_cli
	$(INSTALL) -D $(BUILD_DIR)/gnoi_client $(DESTDIR)/usr/sbin/gnoi_client

deinstall:
	rm $(DESTDIR)/usr/sbin/telemetry
	rm $(DESTDIR)/usr/sbin/dialout_client_cli
	rm $(DESTDIR)/usr/sbin/gnmi_get
	rm $(DESTDIR)/usr/sbin/gnmi_set
	rm $(DESTDIR)/usr/sbin/gnoi_client


