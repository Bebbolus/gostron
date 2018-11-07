# Makefile
#export GOPATH := $(shell pwd)

all:
	go build -buildmode=plugin -o plugins/handlers/pippoHandler.so plugins/handlers/pippoHandler.go
	go build -buildmode=plugin -o plugins/middlewares/Method.so plugins/middlewares/Method.go
	echo $$GOPATH
	#go get -d
	go run *.go


build:
	echo $$GOPATH
	#go get -d
	go build -o out.bin
