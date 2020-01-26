WORK_DIR   := $(shell pwd)

export GITEA_SDK_TEST_URL ?= http://localhost:3000
export GITEA_SDK_TEST_USERNAME ?= test01
export GITEA_SDK_TEST_PASSWORD ?= test01

.PHONY: all
all: clean test build

.PHONY: help
help:
	@echo "Make Routines:"
	@echo " - \"\"              run \"make clean test build\""
	@echo " - build             build sdk"
	@echo " - clean             clean"
	@echo " - fmt               format the code"
	@echo " - lint              run golint"
	@echo " - vet               examines Go source code and reports"
	@echo " - test              run unit tests (need a running gitea)"
	@echo " - test-instance     start a gitea instance for test"


.PHONY: clean
clean:
	rm -r -f test
	go clean -i ./...

.PHONY: fmt
fmt:
	find . -name "*.go" -type f ! -path "./vendor/*" ! -path "./benchmark/*" | xargs gofmt -s -w

.PHONY: vet
vet:
	cd gitea && go vet ./...

.PHONY: lint
lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u golang.org/x/lint/golint; \
	fi
	cd gitea && golint -set_exit_status

.PHONY: test
test:
	@if [ -z "$(shell curl --noproxy "*" "${GITEA_SDK_TEST_URL}/api/v1/version" 2> /dev/null)" ]; then \echo "No test-instance detected!"; exit 1; else \
	cd gitea && go test -cover -coverprofile coverage.out; \
	fi

.PHONY: test-instance
test-instance:
	rm -r ${WORK_DIR}/test 2> /dev/null; \
	mkdir -p ${WORK_DIR}/test/conf/ ${WORK_DIR}/test/data/
	wget "https://dl.gitea.io/gitea/master/gitea-master-linux-amd64" -O ${WORK_DIR}/test/gitea-master; \
	chmod +x ${WORK_DIR}/test/gitea-master; \
	echo "[security]" > ${WORK_DIR}/test/conf/app.ini; \
	echo "INTERNAL_TOKEN = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYmYiOjE1NTg4MzY4ODB9.LoKQyK5TN_0kMJFVHWUW0uDAyoGjDP6Mkup4ps2VJN4" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "INSTALL_LOCK   = true" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "SECRET_KEY     = 2crAW4UANgvLipDS6U5obRcFosjSJHQANll6MNfX7P0G3se3fKcCwwK3szPyGcbo" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "PASSWORD_COMPLEXITY = off" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "[database]" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "DB_TYPE = sqlite3" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "[repository]" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "ROOT = ${WORK_DIR}/test/data/" >> ${WORK_DIR}/test/conf/app.ini; \
	${WORK_DIR}/test/gitea-master migrate -c ${WORK_DIR}/test/conf/app.ini; \
	${WORK_DIR}/test/gitea-master admin create-user --username=${GITEA_SDK_TEST_USERNAME} --password=${GITEA_SDK_TEST_PASSWORD} --email=test01@gitea.io --admin=true --must-change-password=false --access-token -c ${WORK_DIR}/test/conf/app.ini; \
	${WORK_DIR}/test/gitea-master web -c ${WORK_DIR}/test/conf/app.ini

.PHONY: bench
bench:
	cd gitea && go test -run=XXXXXX -benchtime=10s -bench=. || exit 1

.PHONY: build
build:
	cd gitea && go build
