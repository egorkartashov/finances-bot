CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
MOCKGEN=${BINDIR}/mockgen_${GOVER}
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=gitlab.ozon.dev/egor.linkinked/kartashov-egor/cmd/bot

export COMPOSE_PROJECT_NAME=finances-bot

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

all: format build test lint

build: bindir
	go build -o ${BINDIR}/bot ${PACKAGE}

test:
	go test ./...

dev:
	go run ${PACKAGE} -devmode

prod:
	go run ${PACKAGE} 2>&1 | tee ./data/log.txt

generate: install-mockgen
	${MOCKGEN} -source=internal/model/messages/incoming_msg.go -destination=internal/mocks/messages/messages_mocks.go

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-mockgen: bindir
	test -f ${MOCKGEN} || \
		(GOBIN=${BINDIR} go install github.com/golang/mock/mockgen@v1.6.0 && \
		mv ${BINDIR}/mockgen ${MOCKGEN})

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

.PHONY: storage
storage:
	cd deployments && make storage

migrate-up:
	goose --dir migrations postgres ${GOOSE_DB_DSN} up

migrate-down:
	goose --dir migrations postgres ${GOOSE_DB_DSN} down

.PHONY: logs
logs:
	cd deployments && make logs

.PHONY: tracing
tracing:
	cd deployments && make tracing

.PHONY: metrics
metrics:
	cd deployments && make metrics

.PHONY: eventbus
eventbus:
	cd deployments && make eventbus

all: logs tracing metrics

