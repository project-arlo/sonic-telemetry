all: precheck deps telemetry
GO=/usr/local/go/bin/go

TOP_DIR := $(abspath ..)
TELEM_DIR := $(abspath .)
GOFLAGS:=
BUILD_DIR=build
GO_DEP_PATH=$(abspath .)/$(BUILD_DIR)
GO_MGMT_PATH=$(TOP_DIR)/sonic-mgmt-framework
GO_SONIC_TELEMETRY_PATH=$(TOP_DIR)
CVL_GOPATH=$(GO_MGMT_PATH)/build/gopkgs:$(GO_MGMT_PATH):$(GO_MGMT_PATH)/src/cvl/build
GOPATH = $(CVL_GOPATH):$(GO_DEP_PATH):$(GO_MGMT_PATH):/tmp/go:$(GO_SONIC_TELEMETRY_PATH):$(TELEM_DIR)
INSTALL := /usr/bin/install

SRC_FILES=$(shell find . -name '*.go' | grep -v '_test.go' | grep -v '/tests/')
TEST_FILES=$(wildcard *_test.go)
TELEMETRY_TEST_DIR = $(GO_MGMT_PATH)/build/tests/gnmi_server
TELEMETRY_TEST_BIN = $(TELEMETRY_TEST_DIR)/server.test

.PHONY : all precheck deps telemetry clean cleanall check install deinstall

ifdef DEBUG
	GOFLAGS += -gcflags="all=-N -l"
endif

all: deps telemetry $(TELEMETRY_TEST_BIN)

precheck:
	$(shell mkdir -p $(BUILD_DIR))

deps: $(BUILD_DIR)/.deps

$(BUILD_DIR)/.deps: $(MAKEFILE_LIST)
	GOPATH=$(GO_DEP_PATH) $(GO) get -u  github.com/Workiva/go-datastructures/queue
	GOPATH=$(GO_DEP_PATH) $(GO) get -u github.com/openconfig/goyang
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/openconfig/ygot/ygot
	GOPATH=$(GO_DEP_PATH) $(GO) get -u github.com/golang/glog
	GOPATH=$(GO_DEP_PATH) $(GO) get -d go.opentelemetry.io/otel
	cd $(GO_DEP_PATH)/src/go.opentelemetry.io/otel && \
		git clean -fd && git checkout -fq 1f2eba2cdb0fb7e92ce199c6c06167d4080a55ff
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/go-redis/redis
	GOPATH=$(GO_DEP_PATH) $(GO) get -u  github.com/c9s/goprocinfo/linux
	GOPATH=$(GO_DEP_PATH) $(GO) get -u  github.com/golang/protobuf/proto
	GOPATH=$(GO_DEP_PATH) $(GO) get -d  github.com/openconfig/gnmi/proto/gnmi
	GOPATH=$(GO_DEP_PATH) $(GO) get -u  golang.org/x/net/context
	GOPATH=$(GO_DEP_PATH) $(GO) get -u  google.golang.org/grpc
	GOPATH=$(GO_DEP_PATH) $(GO) get -u google.golang.org/grpc/credentials
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/google/gnxi/utils
	cd $(GO_DEP_PATH)/src/github.com/google/gnxi/utils; \
		git reset --hard HEAD; git clean -f -d; git checkout 97477153283d35eb1c7a1b808ebe75bee055e7d8 2>/dev/null; \
		GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/google/gnxi/utils
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

	cd $(GO_DEP_PATH)/src/github.com/openconfig/gnmi/proto/gnmi/; git reset --hard HEAD;git clean -f -d;git checkout e7106f7f5493a9fa152d28ab314f2cc734244ed8 2>/dev/null; true; \
	GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/openconfig/gnmi/proto/gnmi
	cd $(GO_DEP_PATH)/src/github.com/openconfig/ygot/; git reset --hard HEAD;git clean -f -d;git checkout 724a6b18a9224343ef04fe49199dfb6020ce132a 2>/dev/null ; true; \
GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/openconfig/ygot/ygot
	GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/jipanyang/gnxi/utils/xpath
	cd $(GO_DEP_PATH)/src/github.com/jipanyang/gnmi/client/gnmi; git reset --hard HEAD;git clean -f -d;git checkout cb4d464fa018c29eadab98281448000bee4dcc3d 2>/dev/null ; true; \
GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/jipanyang/gnmi/client/gnmi
	
	touch $@

telemetry:$(BUILD_DIR)/telemetry $(BUILD_DIR)/dialout_client_cli $(BUILD_DIR)/gnmi_get $(BUILD_DIR)/gnmi_set $(BUILD_DIR)/gnmi_cli $(BUILD_DIR)/gnoi_client

$(BUILD_DIR)/telemetry:src/telemetry/telemetry.go
	@echo "Building $@"
	GOPATH=$(GOPATH) $(GO) build $(GOFLAGS) -o $@ $^
$(BUILD_DIR)/dialout_client_cli:src/dialout/dialout_client_cli/dialout_client_cli.go
	GOPATH=$(GOPATH) $(GO) build $(GOFLAGS) -o $@ $^
$(BUILD_DIR)/gnmi_get:src/gnmi_clients/gnmi_get.go
	GOPATH=$(GOPATH) $(GO) build $(GOFLAGS) -o $@ $^
$(BUILD_DIR)/gnmi_set:src/gnmi_clients/gnmi_set.go
	GOPATH=$(GOPATH) $(GO) build $(GOFLAGS) -o $@ $^
$(BUILD_DIR)/gnmi_cli:src/gnmi_clients/src/github.com/openconfig/gnmi
	GOPATH=$(PWD)/src/gnmi_clients:$(GOPATH) $(GO) build $(GOFLAGS) -o $@ github.com/openconfig/gnmi/cmd/gnmi_cli
$(BUILD_DIR)/gnoi_client:src/gnmi_clients/gnoi_client.go
	GOPATH=$(PWD)/src/gnmi_clients:$(GOPATH) $(GO) build $(GOFLAGS) -o $@ $^

clean:
	rm -rf $(BUILD_DIR)/telemetry
	rm -rf $(TELEMETRY_TEST_DIR)

cleanall:
	rm -rf $(BUILD_DIR)
	rm -rf $(TELEMETRY_TEST_DIR)

check:
	-$(GO) test -v ${GOPATH}/src/gnmi_server

$(TELEMETRY_TEST_BIN): $(TEST_FILES) $(SRC_FILES)
	GOPATH=$(GOPATH) $(GO) test -c -cover gnmi_server -o $@
	cp -r src/testdata $(TELEMETRY_TEST_DIR)
	cp test/01_create_MyACL1_MyACL2.json $(TELEMETRY_TEST_DIR)
	cp -r $(GO_MGMT_PATH)/debian/sonic-mgmt-framework/usr/sbin/schema $(TELEMETRY_TEST_DIR)


install:
	$(INSTALL) -D $(BUILD_DIR)/telemetry $(DESTDIR)/usr/sbin/telemetry
	$(INSTALL) -D $(BUILD_DIR)/dialout_client_cli $(DESTDIR)/usr/sbin/dialout_client_cli
	$(INSTALL) -D $(BUILD_DIR)/gnmi_get $(DESTDIR)/usr/sbin/gnmi_get
	$(INSTALL) -D $(BUILD_DIR)/gnmi_set $(DESTDIR)/usr/sbin/gnmi_set
	$(INSTALL) -D $(BUILD_DIR)/gnmi_cli $(DESTDIR)/usr/sbin/gnmi_cli
	$(INSTALL) -D $(BUILD_DIR)/gnoi_client $(DESTDIR)/usr/sbin/gnoi_client

	mkdir -p $(DESTDIR)/usr/bin/
	cp -r $(GO_MGMT_PATH)/debian/sonic-mgmt-framework/usr/sbin/schema $(DESTDIR)/usr/sbin
	cp -r $(GO_MGMT_PATH)/debian/sonic-mgmt-framework/usr/sbin/schema $(DESTDIR)/usr/bin
	cp -r $(GO_MGMT_PATH)/debian/sonic-mgmt-framework/usr/models $(DESTDIR)/usr/

deinstall:
	rm $(DESTDIR)/usr/sbin/telemetry
	rm $(DESTDIR)/usr/sbin/dialout_client_cli
	rm $(DESTDIR)/usr/sbin/gnmi_get
	rm $(DESTDIR)/usr/sbin/gnmi_set
	rm $(DESTDIR)/usr/sbin/gnoi_client
