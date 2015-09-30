DESTDIR := /usr/local
DESTBIN := $(DESTDIR)/bin
DESTMAN := $(DESTDIR)/share/man/man1

REPO_NAME := git-comment
REGISTRY_OWNER := stuartnelson3

INSTALLCMD    := install -C
INSTALLDIRCMD := install -d

PROJECT=.
PACKAGES=exec log git search
VERSION=$(shell cat VERSION)
BIN_FILES=git-comment git-comment-grep git-comment-log git-comment-remote git-comment-web
SRC_FILES=comment.go diff.go errors.go file_ref.go lookup.go person.go property_blob.go \
					remote.go storage.go version.go \
					exec/editor.go exec/exec.go exec/pager.go exec/term.go \
					git/commit.go git/commit_range.go git/config.go git/remote.go git/repo.go \
					git/result.go git/var.go \
					log/diff_printer.go log/formatter.go \
					search/formatter.go search/printer.go search/search.go
TEST_FILES=comment_test.go file_ref_test.go person_test.go property_blob_test.go \
					 storage_test.go version_test.go exec/term_test.go git/remote_test.go \
					 log/formatter_test.go
MANSRC=docs/man
MAN_FILES=$(foreach bin,$(BIN_FILES),$(bin).pod)

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
GO15VENDOREXPERIMENT=1

BUILD_DIR=build
BUILD_BIN_DIR=$(BUILD_DIR)/bin
BUILD_MAN_DIR=$(BUILD_DIR)/man
BUILD_BIN_FILES=$(foreach bin,$(BIN_FILES),$(BUILD_BIN_DIR)/$(bin))
BUILD_MAN_FILES=$(foreach bin,$(BIN_FILES),$(BUILD_MAN_DIR)/$(bin).1)

MAN_TITLE=Git Comment Manual
MAN_CMD=pod2man --center="$(MAN_TITLE)" --release="$(VERSION)"

default: build

build: $(BUILD_BIN_FILES)

$(BUILD_BIN_DIR)/%: $(GOPATHSRC_FILES) $(GOPATHPKG_DEPS) bin/%.go
	@$(INSTALLDIRCMD) $(BUILD_BIN_DIR)
	go build \
		-ldflags "-X main.buildVersion=$(VERSION)" \
		-o $(BUILD_BIN_DIR)/$* bin/$*.go

ci: build test

deploy_website:
	git checkout -B gh-pages
	git filter-branch -f --subdirectory-filter docs/git-comment.com
	git clean -df
	# git push -f origin gh-pages

$(BUILD_MAN_DIR)/%.1: $(MANSRC)/%.pod
	@$(INSTALLDIRCMD) $(BUILD_MAN_DIR)
	$(MAN_CMD) $(MANSRC)/$*.pod > $(BUILD_MAN_DIR)/$*.1
	chmod 444 $(BUILD_MAN_DIR)/$*.1

doc: $(BUILD_MAN_FILES)

$(DESTBIN):
	@$(INSTALLDIRCMD) $(DESTBIN)

$(DESTMAN):
	@$(INSTALLDIRCMD) $(DESTMAN)

install: $(DESTBIN) $(DESTMAN) $(BUILD_BIN_FILES) $(BUILD_MAN_FILES)
	@$(INSTALLCMD) $(BUILD_BIN_FILES) $(DESTBIN)
	@$(INSTALLCMD) $(BUILD_MAN_FILES) $(DESTMAN)
	@echo Successfully installed git-comment.

uninstall:
	rm $(foreach bin,$(BIN_FILES), $(DESTMAN)/$(bin).1 $(DESTBIN)/$(bin));

test:
	go test $(foreach pkg,$(PACKAGES),$(PROJECT)/$(pkg)/...)

build-docker:
	docker build --no-cache --force-rm -t $(REGISTRY_OWNER)/$(REPO_NAME):latest .

push-docker: build-docker
	docker push $(REGISTRY_OWNER)/$(REPO_NAME)

test-docker:
	docker run -w /go/src/github.com/git-comment/$(REPO_NAME) -v $(shell pwd):/go/src/github.com/git-comment/$(REPO_NAME) $(REGISTRY_OWNER)/$(REPO_NAME):latest bash -c "go test $(foreach pkg,$(PACKAGES),$(PROJECT)/$(pkg)/...)"

