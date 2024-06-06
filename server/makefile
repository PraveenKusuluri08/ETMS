GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
GOPATH=$(shell go env GOPATH)

default: swag-fmt swag-init build run

swag-fmt:
	swag fmt

swag-init:
	swag init

build:
	go build -o etms-server main.go

run:
	go run main.go

