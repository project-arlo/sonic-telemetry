module github.com/Azure/sonic-telemetry

go 1.12

require (
	github.com/Azure/sonic-mgmt-framework v0.0.0-00010101000000-000000000000
	github.com/Workiva/go-datastructures v1.0.52
	github.com/c9s/goprocinfo v0.0.0-20191125144613-4acdd056c72d
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/golang/protobuf v1.4.0-rc.4
	github.com/google/gnxi v0.0.0-20191016182648-6697a080bc2d
	github.com/jipanyang/gnmi v0.0.0-20180820232453-cb4d464fa018
	github.com/jipanyang/gnxi v0.0.0-20181221084354-f0a90cca6fd0
	github.com/kylelemons/godebug v1.1.0
	github.com/msteinert/pam v0.0.0-20190215180659-f29b9f28d6f9
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/openconfig/gnmi v0.0.0-20190823184014-89b2bf29312c
	github.com/openconfig/gnoi v0.0.0-20191206155121-b4d663a26026
	github.com/openconfig/ygot v0.6.1-0.20190723223108-724a6b18a922
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/tinylib/msgp v1.1.1 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	gonum.org/v1/plot v0.0.0-20200226011204-b25252b0d522 // indirect
	google.golang.org/grpc v1.27.1
)

replace github.com/Azure/sonic-mgmt-framework => ../sonic-mgmt-framework
