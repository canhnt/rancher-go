PKG = github.com/canhnt/rancher-go/cmd/rancherctl

# TODO Make it DRY
.PHONY: build_mac
build_mac:
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -o rancherctl_darwin_amd64 ${PKG}

.PHONY: build_linux
build_linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o rancherctl_linux_amd64 ${PKG}

.PHONY: build
build: build_mac build_linux

.PHONY: clean
clean:
		rm -rf rancherctl_darwin_amd64 rancherctl_linux_amd64
