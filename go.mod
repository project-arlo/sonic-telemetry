module github.com/Azure/sonic-telemetry

go 1.12

require (
	github.com/Workiva/go-datastructures v1.0.50
	github.com/antchfx/jsonquery v1.1.0
	github.com/antchfx/xmlquery v1.1.1-0.20191015122529-fe009d4cc63c
	github.com/antchfx/xpath v1.1.2
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/c9s/goprocinfo v0.0.0-20191125144613-4acdd056c72d
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/facette/natsort v0.0.0-20181210072756-2cd4dd1e2dcb
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/godbus/dbus v4.1.0+incompatible
	github.com/gogo/protobuf v1.3.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/google/gnxi v0.0.0-20191016182648-6697a080bc2d
	github.com/jipanyang/gnmi v0.0.0-20180820232453-cb4d464fa018
	github.com/jipanyang/gnxi v0.0.0-20181221084354-f0a90cca6fd0
	github.com/kylelemons/godebug v1.1.0
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/msteinert/pam v0.0.0-20190215180659-f29b9f28d6f9
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/openconfig/gnmi v0.0.0-20190823184014-89b2bf29312c
	github.com/openconfig/gnoi v0.0.0-20191206155121-b4d663a26026
	github.com/openconfig/goyang v0.0.0-20190924211109-064f9690516f
	github.com/openconfig/ygot v0.6.1-0.20190723223108-724a6b18a922
	github.com/pborman/getopt v0.0.0-20190409184431-ee0cd42419d3 // indirect
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/philopon/go-toposort v0.0.0-20170620085441-9be86dbd762f
	github.com/pkg/profile v1.4.0
	github.com/tinylib/msgp v1.1.1 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
	golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3 // indirect
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	golang.org/x/text v0.3.2
	golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135 // indirect
	gonum.org/v1/plot v0.0.0-20200226011204-b25252b0d522 // indirect
	google.golang.org/grpc v1.25.1
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc // indirect
)

replace github.com/Azure/sonic-mgmt-framework => ../sonic-mgmt-framework
