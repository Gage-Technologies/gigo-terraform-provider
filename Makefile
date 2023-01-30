default: testacc

fmt:
	terraform fmt -recursive

gen:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

build:
	go build -o bin/gigo-provider

build-dev:
	make build
	mkdir -p ~/.terraform.d/plugins/gigo.dev/gigo/gigo-dev/lastest/
	cp bin/gigo-provider ~/.terraform.d/plugins/gigo.dev/gigo/gigo-dev/lastest/linux_amd64

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m