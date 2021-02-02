pkgname = kamailio-jsonrpc-client
namespace = voipxswitch
modname = github.com/$(namespace)/$(pkgname)

default: clean init build

init:
	go mod init $(modname)

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(pkgname)

clean:
	rm -rf go.mod
	rm -rf go.sum
