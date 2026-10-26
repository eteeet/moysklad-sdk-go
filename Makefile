FIXTURES_REF ?= HEAD

.PHONY: fixtures fixtures-commit example

fixtures:
	go run ./testdata/syncfixtures -ref "$(FIXTURES_REF)"

fixtures-commit:
	go run ./testdata/syncfixtures -ref "$(FIXTURES_REF)" -commit

example:
	@set -ea; \
	. ./.env; \
	go run ./examples/product_get
