export GOPATH=$(shell pwd)/../

all:
	go get -u github.com/endurox-dev/endurox-go
	go get -u github.com/golang/protobuf/protoc-gen-go
	$(MAKE) -C userdet
	$(MAKE) -C ubftab
	$(MAKE) -C userclt
	$(MAKE) -C userclt_ubf
	$(MAKE) -C usersrv
	$(MAKE) -C usersrv_ubf
	
clean:
	- $(MAKE) -C userdet clean
	- $(MAKE) -C ubftab clean
	- $(MAKE) -C userclt clean
	- $(MAKE) -C userclt_ubf clean
	- $(MAKE) -C usersrv clean
	- $(MAKE) -C usersrv_ubf clean
	- go get -u github.com/endurox-dev/endurox-go
	- go get -u github.com/golang/protobuf/protoc-gen-go

.PHONY: clean all
