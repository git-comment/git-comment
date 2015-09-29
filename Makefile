DESTDIR := /usr/local
DESTBIN := $(DESTDIR)/bin
DESTMAN := $(DESTDIR)/share/man/man1

INSTALLCMD    := install -C
INSTALLDIRCMD := install -d

PROJECT=libgitcomment
PACKAGES=exec log git search
VERSION=$(shell cat VERSION)
DEPENDENCIES=gopkg.in/libgit2/git2go.v23 \
	github.com/stvp/assert \
	github.com/cevaris/ordered_map \
  gopkg.in/alecthomas/kingpin.v2 \
  github.com/kylef/result.go/src/result \
  github.com/blevesearch/bleve \
  github.com/blang/semver
BIN_FILES=$(basename $(shell ls bin))
SRC_FILES=$(filter-out test,$(shell git ls-files "$(PROJECT)/**.go"))
TEST_FILES=$(shell git ls-files "$(PROJECT)/**_test.go")
MANSRC=docs/man
MAN_FILES=$(foreach bin,$(BIN_FILES),$(bin).pod)

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
GOPATH=$(shell pwd)/_workspace
GO=GOPATH="$(GOPATH)" go
GOPATHSRC=$(GOPATH)/src/$(PROJECT)
GOPATHSRC_FILES=$(addprefix $(GOPATH)/src/,$(SRC_FILES))
GOPATHSRC_TESTS=$(addprefix $(GOPATH)/src/,$(TEST_FILES))
GOPATHPKG=$(GOPATH)/pkg/$(GOOS)_$(GOARCH)
GOPATHPKG_DEPS=$(foreach dep,$(DEPENDENCIES),$(GOPATHPKG)/$(dep).a)

BUILD_DIR=build
BUILD_BIN_DIR=$(BUILD_DIR)/bin
BUILD_MAN_DIR=$(BUILD_DIR)/man
BUILD_BIN_FILES=$(foreach bin,$(BIN_FILES),$(BUILD_BIN_DIR)/$(bin))
BUILD_MAN_FILES=$(foreach bin,$(BIN_FILES),$(BUILD_MAN_DIR)/$(bin).1)

MAN_TITLE=Git Comment Manual
MAN_CMD=pod2man --center="$(MAN_TITLE)" --release="$(VERSION)"

default: build

build: $(GOPATHSRC_FILES) $(BUILD_BIN_FILES)

$(BUILD_BIN_DIR)/%: $(GOPATHSRC_FILES) $(GOPATHPKG_DEPS) bin/%.go
	@$(INSTALLDIRCMD) $(BUILD_BIN_DIR)
	$(GO) build -ldflags "-X main.buildVersion=$(VERSION)" -o $(BUILD_BIN_DIR)/$* bin/$*.go

$(GOPATHSRC)/%.go: $(GOPATHPKG_DEPS) $(PROJECT)/%.go
	@$(INSTALLDIRCMD) $(GOPATHSRC)/$(dir $*)
	@$(INSTALLCMD) $(PROJECT)/$*.go $(GOPATHSRC)/$*.go

$(GOPATHPKG)/%.a:
	$(GO) get $*
	@rm -rf $(GOPATH)/src/$*/.git

ci: build test

clean:
	$(GO) clean $(PROJECT) || true
	rm -rf $(GOPATHSRC) $(BUILD_DIR)

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

test: $(GOPATHSRC_FILES) $(GOPATHSRC_TESTS)
	$(GO) test $(PROJECT) $(foreach pkg,$(PACKAGES),$(PROJECT)/$(pkg));
