.PHONY: gox

default: xcompile

gox:
	go get -v github.com/mitchellh/gox
	gox -verbose -build-toolchain -osarch="darwin/amd64 linux/amd64"

xcompile: client/main.go server/main.go
	gox -verbose -osarch="darwin/amd64 linux/amd64" -output="build/client-{{.OS}}_{{.Arch}}" github.com/fgrehm/notify-send-http/client
	gox -verbose -osarch="darwin/amd64 linux/amd64" -output="build/server-{{.OS}}_{{.Arch}}" github.com/fgrehm/notify-send-http/server
