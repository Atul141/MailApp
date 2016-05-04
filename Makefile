# Please refer to http://clarkgrubb.com/makefile-style-guide
#
ENV_MIN_COVERAGE = 46

conf_file := default.conf
mailbox_root := $(shell pwd | grep -o .*/git\.mailbox\.com)

# Prologue

MAKEFLAGS += --warn-undefined-variables
SHELL := bash
# SHELLFLAGS has no effect on GNU Make < 3.82
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:

# Environment
ifndef ENV_DO_LINT
ENV_DO_LINT := 1
endif

ifndef ENV_MIN_COVERAGE
ENV_MIN_COVERAGE := 85
endif

ifndef ENV_BUILD_LINUX
ENV_BUILD_LINUX := 1
endif

ifeq ($(origin GOPATH), undefined)
$(error GOPATH environment variable not set!)
endif

# TODO: Needs to be removed in Go 1.6 / 1.7
GO15VENDOREXPERIMENT := 1
export GO15VENDOREXPERIMENT

# common variables
os := $(shell go env GOOS)
package_basename := $(shell pwd | grep -o git\.mailbox\.com/.*)
apps_folder = apps

out_test_path := out/test
out_build_path := out/build

main_paths := $(shell find ./$(apps_folder) -mindepth 1 -maxdepth 1 -type d)

binary_names := $(main_paths:./$(apps_folder)/%=$(out_build_path)/%)
binaries :=

ifeq ($(os), linux)
	linux_binaries = $(addsuffix -linux,$(binary_names))
	binaries += $(linux_binaries)
endif
ifeq ($(os), windows)
	local_binaries = $(addsuffix -windows.exe,$(binary_names))
	binaries += $(local_binaries)
	ifeq ($(ENV_BUILD_LINUX), 1)
		linux_binaries := $(addsuffix -linux,$(binary_names))
		binaries += $(linux_binaries)
	endif
endif
ifeq ($(os), darwin)
	local_binaries = $(addsuffix -darwin,$(binary_names))
	binaries += $(local_binaries)
	ifeq ($(ENV_BUILD_LINUX), 1)
		linux_binaries = $(addsuffix -linux,$(binary_names))
		binaries += $(linux_binaries)
	endif
endif

sources := $(wildcard **/*.go)
version := $(shell git describe --tags --always)
ldflags := -ldflags "-X main.version=$(version)"

# NOTE that exported variables are available to recipes, but not to shell commands
# Hence passing GO15VENDOREXPERIMENT explicitly to glide
top_package_patterns := $(subst ./gen/...,,$(shell GO15VENDOREXPERIMENT=1 glide nv))
top_package_names = $(top_package_patterns:./%/...=%)
all_package_paths := $(shell GOPATH=$(GOPATH) go list ${top_package_patterns})
all_package_names = $(all_package_paths:$(package_basename)/%=%)
all_package_tests = $(all_package_names:%=test.%)
all_package_test_reports = $(all_package_names:%=$(out_test_path)/%.func.txt)
top_package_tests = $(top_package_names:%=test.%)
top_package_test_reports = $(top_package_names:%=$(out_test_path)/%.func.txt)
all_test_report = $(out_test_path)/all.func.txt
all_test_report_html = $(out_test_path)/all.func.html
minimum_coverage_percent := $(ENV_MIN_COVERAGE)
db_create_script_path := "./db/db_create.sql"
db_destroy_script_path := "./db/db_destroy.sql"
db_seed_script_path := "./db/db_seed.sql"

.PHONY: info
info:
	@echo "Build information"
	@go version
	@echo  "OS           :: $(os)"
	@echo  "GOPATH       :: $(GOPATH)"

.PHONY: setup
setup:
	go get -v github.com/golang/lint/golint
	./setup_glide.sh
	go get bitbucket.org/liamstask/goose/cmd/goose

.PHONY: glide.update
glide.update:
	glide update

.PHONY: glide.install
glide.install:
	glide install

# compile
$(out_build_path)/%-windows.exe: $(sources)
	@echo "Building windows executable $@"
	@GOOS=windows GOARCH=amd64 go build $(ldflags) -o $@ ./$(apps_folder)/$*

$(out_build_path)/%-darwin: $(sources)
	@echo "Building darwin executable $@"
	@GOOS=darwin GOARCH=amd64 go build $(ldflags) -o $@ ./$(apps_folder)/$*

$(out_build_path)/%-linux: $(sources)
	@echo "Building linux executable $@"
	@GOOS=linux GOARCH=amd64 go build $(ldflags) -o $@ ./$(apps_folder)/$*

.PHONY: compile
compile: $(binaries)

.PHONY: clean.compile
clean.compile:
	@echo "Removing binaries, if any"
	@rm -f $(binaries)

.PHONY: clean.gen
clean.gen:
	@echo "Removing generated sources, if any"
	@rm -rf gen/

.PHONY: install
install:
	@go install $(ldflags) $(main_paths)

# unit+integration test
$(out_test_path)/%.func.txt: $(sources)
	@mkdir -p $(@D)/$*
	@echo "Testing top level package $*";
	@result=$$( find $* -name '*.go' -print0 | xargs -0 -n 1 dirname | uniq 2>&1 ); \
	for package in $$result; do \
		echo "Testing package $$package"; \
		profilefile=$(out_test_path)/$$package; mkdir -p `dirname $$profilefile`; \
		go test -covermode="count" -coverprofile="$(out_test_path)/$$package.func.txt" ./$$package; \
	done;
	@find $(@D)/$* -name '*.func.txt' | xargs cat >> $@
	@rm -rf $(@D)/$*


$(all_test_report): $(top_package_test_reports)
	@echo "Creating consolidated coverage file $@"
	@{ echo "mode: count"; \
	  cat $^ | \
	  sed '/^mode.*count$$/d' | sed '/^$$/d' | sed 's/\r$$/$$/' ; } > $@.tmp
	@mv $@.tmp $@

$(all_test_report_html): $(all_test_report)
	@go tool cover --html $< -o $@

.PHONY: $(top_package_tests)
$(top_package_tests):test.%:$(out_test_path)/%.func.txt

.PHONY: test
test: coverage = $(shell GOPATH=$(GOPATH) go tool cover --func=out/test/all.func.txt | tail -1 | awk '{ print int($$3) }' | sed 's/%$$//')
test: $(all_test_report) $(all_test_report_html) 
	$(info Total Coverage = $(coverage)%)
	@if [[ $(coverage) -lt $(minimum_coverage_percent) ]]; then \
		echo "Coverage ${coverage} is below $(minimum_coverage_percent)%! Failing build." ;\
		exit 1 ;\
	fi

.PHONY: clean.test
clean.test:
	@echo "Removing test reports, if any"
	@rm -rf $(out_test_path)

# check code format and style
.PHONY: fmt
fmt:
	@echo "Checking formatting of go sources"
	@result=$$(gofmt -d -l -e $(top_package_names) 2>&1); \
		if [[ "$$result" ]]; then \
			echo "$$result"; \
			echo 'gofmt failed!'; \
			exit 1; \
		fi

# Database
.PHONY: db.create
db.create:
	@if [[ -d "db" ]]; then \
		echo "Creating database"; \
		psql -f $(db_create_script_path); \
	fi

.PHONY: db.migrate
db.migrate:
	@if [[ -d "db" ]]; then \
		echo "Migrating database"; \
		goose up; \
	fi

.PHONY: db.rollback
db.rollback:
	@if [[ -d "db" ]]; then \
		echo "Rolling back database"; \
		goose down; \
	fi


.PHONY: db.seed
db.seed:
	@if [[ -d "db" ]]; then \
		echo "Seeding database"; \
		psql -f $(db_seed_script_path); \
	fi

.PHONY: db.destroy
db.destroy:
	@if [[ -d "db" ]]; then \
		echo "Destroying database"; \
		psql -f $(db_destroy_script_path); \
	fi

.PHONY: db.reset
db.reset: db.destroy db.create db.migrate db.seed

# Fix code format and style
# NOT TO BE RUN ON BUILD
.PHONY: fixfmt
fixfmt:
	@echo "Fixing format of go sources"
	@gofmt -w -l -e $(top_package_names) 2>&1; \
		if [[ "$$?" != 0 ]]; then \
		    echo "gofmt failed! (exit-code: '$$?')"; \
		    exit 1; \
		fi

.PHONY: vet
vet:
	@echo "Running go vet"
	@go vet $(all_package_paths)

targetDir := $(addprefix $(GOPATH), /src/$(package_basename))

.PHONY: lint
lint:
ifeq ($(ENV_DO_LINT),1)
	@echo "Running golint"
	@echo $(all_package_paths) | xargs -n 1 golint -min_confidence 0.8
else
	@echo "Skipping lint"
endif

.PHONY: check
check: fmt vet lint test

# clean
.PHONY: clean
clean: clean.compile clean.test

# build
.PHONY: quick
quick: info compile test

.PHONY: build
build: info clean compile check test

# copy default conf to build dir
.PHONY: copyconf
copyconf:
	@cp $(conf_file) $(out_build_path)

# all
all:
	@$(MAKE) -rpn  : 2>/dev/null | grep '^.PHONY' | cut -d' ' -f2- | xargs -n 1 echo | sort

.DEFAULT_GOAL: all


