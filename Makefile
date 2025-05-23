# Makefile to compile the project for Windows amd64

# Set the target OS and architecture
GOOS := windows
GOARCH := amd64

# Define the output binary name
BINARY_NAME := myproject.exe

# Default target
all: build

# Build the project
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME)

# Clean the build
clean:
	rm -f $(BINARY_NAME)

.PHONY: all build clean
