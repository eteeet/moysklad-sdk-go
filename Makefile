FIXTURES_REF ?= HEAD
# Параллельность пакетов в тестах. HTTP-тесты идут через internal/memnet, без loopback TCP.
CHECK_P ?= 2

.PHONY: check fixtures fixtures-commit example example-create example-update

# Проверки выпуска без сети и токенов: gofmt, vet, race-тесты (включая манифест фикстур и примеры),
# пробелы во всём дереве HEAD и в незакоммиченных изменениях.
check:
	@out="$$(gofmt -l .)"; if [ -n "$$out" ]; then echo "gofmt:"; echo "$$out"; exit 1; fi
	go vet ./... ./testdata
	go test -race -count=1 -p $(CHECK_P) ./...
	git diff --check $$(git hash-object -t tree /dev/null) HEAD
	git diff --check HEAD

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
