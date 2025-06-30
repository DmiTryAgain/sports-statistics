NAME := sportstat
MAIN := main.go
PKG := `go list -mod=vendor -f {{.Dir}} ./...`
GOFLAGS=-mod=vendor

ifeq ($(RACE),1)
	GOFLAGS+=-race
endif

build:
	@CGO_ENABLED=0 go build $(GOFLAGS) -o ${NAME} $(MAIN)

run:
	@echo "Compiling"
	@go run $(GOFLAGS) $(MAIN) -config=config/local.toml -dev

test:
	@go test -v ./...

mod:
	@go mod tidy
	@go mod vendor
	@git add vendor

fmt:
	@golangci-lint fmt

lint:
	@golangci-lint version
	@golangci-lint config verify
	@golangci-lint run

MAPPING := "statistic:statistics"
NS := "NONE"
ENTITY := "NONE"

mfd-xml:
	@mfd-generator xml -c "postgres://postgres:postgres@localhost:5432/sport_statsrv?sslmode=disable" -m ./docs/model/sportStatistics.mfd -n $(MAPPING)
mfd-model:
	@mfd-generator model -m ./docs/model/sportStatistics.mfd -p db -o ./pkg/db
mfd-repo: --check-ns
	@mfd-generator repo -m ./docs/model/sportStatistics.mfd -p db -o ./pkg/db -n $(NS)

--check-ns:
ifeq ($(NS),"NONE")
	$(error "You need to set NS variable before run this command.")
endif

--check-entity:
ifeq ($(ENTITY),"NONE")
	$(error "You need to set ENTITY variable before run this command.")
endif
