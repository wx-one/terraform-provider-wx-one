default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cp main.go main.goorig; sed -i 's#registry.terraform.io/providers/wx-one#hashicorp.com/edu#' main.go; cd tools; go generate ./...; cd ..; cp main.goorig main.go

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate
