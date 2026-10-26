FIXTURES_REF ?= HEAD

.PHONY: fixtures fixtures-commit example example-create example-update

fixtures:
	go run ./testdata -ref "$(FIXTURES_REF)"

fixtures-commit:
	go run ./testdata -ref "$(FIXTURES_REF)" -commit

example:
	@set -ea; \
	. ./.env; \
	go run ./examples/product_get

example-create:
	@set -ea; \
	. ./.env; \
	go run ./examples/product_create

example-update:
	@set -ea; \
	. ./.env; \
	go run ./examples/product_update
