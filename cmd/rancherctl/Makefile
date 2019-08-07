PKG = github.com/canhnt/rancher-go/cmd/rancherctl
BUILD_OUTPUT := rancherctl

.PHONY: build
build: ${BUILD_OUTPUT}
${BUILD_OUTPUT}:
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -o ${BUILD_OUTPUT} ${PKG}

.PHONY: clean
clean:
		rm -rf ${BUILD_OUTPUT}