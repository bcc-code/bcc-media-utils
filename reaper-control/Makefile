GOOS := windows
GOARCH := amd64

BINARY_NAME := reapercontrol.exe

all: build

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

.PHONY: all build clean
