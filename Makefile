default: testacc

fmt:
	terraform fmt -recursive

gen:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

build:
	CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' -o bin/terraform-provider-gigo

build-dev:
	make build
	mkdir -p ~/.terraform.d/plugins/terraform.local/gigo/gigo/0.0.1/linux_amd64/
	cp bin/terraform-provider-gigo ~/.terraform.d/plugins/terraform.local/gigo/gigo/0.0.1/linux_amd64/

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m