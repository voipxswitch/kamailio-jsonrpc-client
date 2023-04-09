pkgname = kamailio-jsonrpc-client
namespace = voipxswitch
modname = github.com/$(namespace)/$(pkgname)
vtag = $(shell git rev-parse --short HEAD)
dtag = $(shell date +%s)-$(vtag)

default: clean init build

init:
	go mod init $(modname)
	go mod tidy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(pkgname)
	@echo tagged release $(vtag)
	$(shell echo "{\"tag\":\"$(vtag)\"}" > tag.json)

clean:
	rm -rf go.mod
	rm -rf go.sum
