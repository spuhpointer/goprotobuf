
SOURCES := $(shell find $(SOURCEDIR) -name '*.proto')
OUTPUT := $(patsubst %.proto,%.pb.go,$(SOURCES))


$(info $$SOURCES is [${SOURCES}])
$(info $$OUTPUT is [${OUTPUT}])

$(OUTPUT): $(SOURCES)
	protoc --go_out=. *.proto
	go build
	go install ./...
	
.DEFAULT_GOAL: $(OUTPUT)

.PHONY: clean

clean:
	- rm *.go

