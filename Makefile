all: precheck deps patch telemetry
GO=/usr/local/go/bin/go

TOP_DIR := $(abspath ..)
TELEM_DIR := $(abspath .)
GOFLAGS:=
BUILD_DIR=build
GO_DEP_PATH=$(abspath .)/$(BUILD_DIR)
GO_MGMT_PATH=$(TOP_DIR)/sonic-mgmt-framework
GO_SONIC_TELEMETRY_PATH=$(TOP_DIR)

GOPATH = $(GO_DEP_PATH):$(GO_MGMT_PATH):/tmp/go:$(GO_SONIC_TELEMETRY_PATH):$(TELEM_DIR)
INSTALL := /usr/bin/install

SRC_FILES=$(shell find . -name '*.go' | grep -v '_test.go' | grep -v '/tests/')
TEST_FILES=$(wildcard *_test.go)
TELEMETRY_TEST_DIR = $(GO_MGMT_PATH)/build/tests/gnmi_server
TELEMETRY_TEST_BIN = $(TELEMETRY_TEST_DIR)/server.test

.PHONY : all precheck deps patch telemetry clean cleanall check install deinstall

ifdef DEBUG
	GOFLAGS += -gcflags="all=-N -l"
endif

all: deps patch telemetry $(TELEMETRY_TEST_BIN)

precheck:
	$(shell mkdir -p $(BUILD_DIR))

deps: $(BUILD_DIR)/.deps

$(BUILD_DIR)/.deps:
	touch $(BUILD_DIR)/.deps
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/openconfig/gnmi; cd $(GO_DEP_PATH)/src/github.com/openconfig/gnmi; \
git checkout 89b2bf29312cda887da916d0f3a32c1624b7935f 2>/dev/null ; true; \
	GOPATH=$(GO_DEP_PATH) $(GO) get -d  github.com/Workiva/go-datastructures/queue; cd $(GO_DEP_PATH)/src/github.com/Workiva/go-datastructures; \
git checkout f07cbe3f82ca2fd6e5ab94afce65fe43319f675f 2>/dev/null ; true; \
# GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/Workiva/go-datastructures
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/openconfig/goyang; cd $(GO_DEP_PATH)/src/github.com/openconfig/goyang; \
git checkout 064f9690516f4f72db189f4690b84622c13b7296 2>/dev/null ; true; \
# GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/openconfig/goyang
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/golang/glog; cd $(GO_DEP_PATH)/src/github.com/golang/glog; \
git checkout 23def4e6c14b4da8ac2ed8007337bc5eb5007998 2>/dev/null ; true; \
# GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/golang/glog
	GOPATH=$(GO_DEP_PATH) $(GO) get -d  github.com/golang/protobuf/proto; cd $(GO_DEP_PATH)/src/github.com/golang/protobuf/proto; \
git checkout ed6926b37a637426117ccab59282c3839528a700 2>/dev/null ; true; \
# GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/golang/protobuf/proto
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/openconfig/ygot; cd $(GO_DEP_PATH)/src/github.com/openconfig/ygot; \
git checkout 724a6b18a9224343ef04fe49199dfb6020ce132a 2>/dev/null ; true; \
# GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/openconfig/ygot/ygot
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/go-redis/redis; cd $(GO_DEP_PATH)/src/github.com/go-redis/redis; \
git checkout d19aba07b47683ef19378c4a4d43959672b7cec8 2>/dev/null ; true; \
# GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/go-redis/redis
	GOPATH=$(GO_DEP_PATH) $(GO) get -d  github.com/c9s/goprocinfo/linux; cd $(GO_DEP_PATH)/src/github.com/c9s/goprocinfo/linux; \
git checkout 0b2ad9ac246b05c4f5750721d0c4d230888cac5e 2>/dev/null ; true; \
# GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/c9s/goprocinfo/linux
	GOPATH=$(GO_DEP_PATH) $(GO) get -d  golang.org/x/net/context
	GOPATH=$(GO_DEP_PATH) $(GO) get -d google.golang.org/grpc
	GOPATH=$(GO_DEP_PATH) $(GO) get -d gopkg.in/go-playground/validator.v9
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/gorilla/mux; cd $(GO_DEP_PATH)/src/github.com/gorilla/mux; \
git checkout 49c01487a141b49f8ffe06277f3dca3ee80a55fa 2>/dev/null ; true; \
# GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/gorilla/mux
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/google/gnxi/utils; cd $(GO_DEP_PATH)/src/github.com/google/gnxi/utils; \
git checkout 6697a080bc2d3287d9614501a3298b3dcfea06df 2>/dev/null ; true; \
# GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/google/gnxi/utils
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/jipanyang/gnxi/utils/xpath; cd $(GO_DEP_PATH)/src/github.com/jipanyang/gnxi/utils/xpath; \
git checkout f0a90cca6fd0041625bcce561b71f849c9b65a8d 2>/dev/null ; true; \
# GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/jipanyang/gnxi/utils/xpath
	GOPATH=$(GO_DEP_PATH) $(GO) get -u github.com/jipanyang/gnmi/client; cd $(GO_DEP_PATH)/github.com/jipanyang/gnmi/client; \
git checkout cb4d464fa018c29eadab98281448000bee4dcc3d 2>/dev/null ; true; \

	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/antchfx/jsonquery; cd $(GO_DEP_PATH)/src/github.com/antchfx/jsonquery; \
git checkout 3535127d6ca5885dbf650204eb08eabf8374a274 2>/dev/null ; true; \
# GOPATH=$(GO_DEP_PATH) $(GO) install -v -gcflags "-N -l" $(GO_DEP_PATH)/src/github.com/antchfx/jsonquery
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/antchfx/xmlquery; cd $(GO_DEP_PATH)/src/github.com/antchfx/xmlquery; \
git checkout 16f1e6cdc5fe44a7f8e2a8c9faf659a1b3a8fd9b 2>/dev/null ; true; \
	GOPATH=$(GO_DEP_PATH) $(GO) get -d golang.org/x/crypto; cd $(GO_DEP_PATH)/src/golang.org/x/crypto; \
git checkout 86a70503ff7e82ffc18c7b0de83db35da4791e6a 2>/dev/null ; true; \
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/xeipuuv/gojsonschema; cd $(GO_DEP_PATH)/src/github.com/xeipuuv/gojsonschema; \
git checkout 001aa27b4d110df2cd68ee5f86f17d8e6d1f3b36 2>/dev/null ; true; \
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/openconfig/gnoi; cd $(GO_DEP_PATH)/src/github.com/openconfig/gnoi; \
git checkout b4d663a260265ce26d9880777b88a6525a1a28d4 2>/dev/null ; true; \
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/msteinert/pam; cd $(GO_DEP_PATH)/src/github.com/msteinert/pam; \
git checkout f29b9f28d6f9a1f6c4e6fd5db731999eb946574b 2>/dev/null ; true; \
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/dgrijalva/jwt-go; cd $(GO_DEP_PATH)/src/github.com/dgrijalva/jwt-go; \
git checkout 5e25c22bd5d6de03265bbe5462dcd162f85046f6 2>/dev/null ; true; \
	GOPATH=$(GO_DEP_PATH) $(GO) get -d gopkg.in/godbus/dbus.v5; cd $(GO_DEP_PATH)/src/gopkg.in/godbus/dbus.v5; \
git checkout 37bf87eef99d69c4f1d3528bd66e3a87dc201472 2>/dev/null ; true; \
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/gogo/protobuf; cd $(GO_DEP_PATH)/src/github.com/gogo/protobuf; \
git checkout 5628607bb4c51c3157aacc3a50f0ab707582b805 2>/dev/null ; true; \
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/philopon/go-toposort; cd $(GO_DEP_PATH)/src/github.com/philopon/go-toposort; \
git checkout 9be86dbd762f98b5b9a4eca110a3f40ef31d0375 2>/dev/null ; true; \
	GOPATH=$(GO_DEP_PATH) $(GO) get -d github.com/facette/natsort; cd $(GO_DEP_PATH)/src/github.com/facette/natsort; \
git checkout 2cd4dd1e2dcba4d85d6d3ead4adf4cfd2b70caf2 2>/dev/null ; true; \

patch: $(BUILD_DIR)/.patched

$(BUILD_DIR)/.patched:
	touch $(BUILD_DIR)/.patched
	patch -p0 <patches/gnmi_cli.all.patch
	patch -p1 -d build/src/github.com/openconfig <$(GO_MGMT_PATH)/ygot-modified-files/ygot.patch
	patch -p1 -d build/src/github.com/openconfig/goyang <$(GO_MGMT_PATH)/goyang-modified-files/goyang.patch
	@grep ParseJsonMap  $(GO_DEP_PATH)/src/github.com/antchfx/jsonquery/node.go || \
	printf "\nfunc ParseJsonMap(jsonMap *map[string]interface{}) (*Node, error) {\n \
		doc := &Node{Type: DocumentNode}\n \
		parseValue(*jsonMap, doc, 1)\n \
		return doc, nil\n \
	}\n" >> $(GO_DEP_PATH)/src/github.com/antchfx/jsonquery/node.go
	touch $@

telemetry:$(BUILD_DIR)/telemetry $(BUILD_DIR)/dialout_client_cli $(BUILD_DIR)/gnmi_get $(BUILD_DIR)/gnmi_set $(BUILD_DIR)/gnmi_cli $(BUILD_DIR)/gnoi_client

$(BUILD_DIR)/telemetry:src/telemetry/telemetry.go
	@echo "Building $@"
	mkdir -p $(GO_MGMT_PATH)/build/cvl/schema
	GOPATH=$(GOPATH) BUILD_GOPATH=$(GO_DEP_PATH) GO=$(GO) $(GO) generate translib/ocbinds
	TOPDIR=$(GO_MGMT_PATH) make -C $(GO_MGMT_PATH)/src/cvl/schema
	make -C $(GO_MGMT_PATH)/models
	make -C $(GO_MGMT_PATH)/models/yang
	make -C $(GO_MGMT_PATH)/models/yang/sonic
	GOPATH=$(GOPATH):$(GO_MGMT_PATH) $(GO) build cvl
	GOPATH=$(GOPATH) $(GO) build $(GOFLAGS) -o $@ $^
$(BUILD_DIR)/dialout_client_cli:src/dialout/dialout_client_cli/dialout_client_cli.go
	GOPATH=$(GOPATH) $(GO) build $(GOFLAGS) -o $@ $^
$(BUILD_DIR)/gnoi_client:src/gnmi_clients/gnoi_client.go
	GOPATH=$(PWD)/src/gnmi_clients:$(GOPATH) $(GO) build $(GOFLAGS) -o $@ $^
$(BUILD_DIR)/gnmi_get:$(BUILD_DIR)/src/github.com/jipanyang/gnxi/gnmi_get/gnmi_get.go
	GOPATH=$(GO_DEP_PATH):$(GOPATH) $(GO) build $(GOFLAGS) -o $@ $^
$(BUILD_DIR)/gnmi_set:$(BUILD_DIR)/src/github.com/jipanyang/gnxi/gnmi_set/gnmi_set.go
	GOPATH=$(GO_DEP_PATH):$(GOPATH) $(GO) build $(GOFLAGS) -o $@ $^
$(BUILD_DIR)/gnmi_cli:$(BUILD_DIR)/src/github.com/openconfig/gnmi
	GOPATH=$(GO_DEP_PATH):$(GOPATH) $(GO) build $(GOFLAGS) -o $@ $^/cmd/gnmi_cli/gnmi_cli.go


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

install:
	$(INSTALL) -D $(BUILD_DIR)/telemetry $(DESTDIR)/usr/sbin/telemetry
	$(INSTALL) -D $(BUILD_DIR)/dialout_client_cli $(DESTDIR)/usr/sbin/dialout_client_cli
	$(INSTALL) -D $(BUILD_DIR)/gnmi_get $(DESTDIR)/usr/sbin/gnmi_get
	$(INSTALL) -D $(BUILD_DIR)/gnmi_set $(DESTDIR)/usr/sbin/gnmi_set
	$(INSTALL) -D $(BUILD_DIR)/gnmi_cli $(DESTDIR)/usr/sbin/gnmi_cli
	$(INSTALL) -D $(BUILD_DIR)/gnoi_client $(DESTDIR)/usr/sbin/gnoi_client
	mkdir -p $(DESTDIR)/usr/bin/
	cp -r $(GO_MGMT_PATH)/src/cvl/schema $(DESTDIR)/usr/sbin
	mkdir -p $(DESTDIR)/usr/models/yang
	find $(GO_MGMT_PATH)/models -name '*.yang' -exec cp {} $(DESTDIR)/usr/models/yang/ \;

deinstall:
	rm $(DESTDIR)/usr/sbin/telemetry
	rm $(DESTDIR)/usr/sbin/dialout_client_cli
	rm $(DESTDIR)/usr/sbin/gnmi_get
	rm $(DESTDIR)/usr/sbin/gnmi_set
	rm $(DESTDIR)/usr/sbin/gnoi_client
