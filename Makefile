DESTDIR := /usr/local

PROJECT=git_comment
PACKAGES=exec log git search
VERSION=$(shell cat VERSION)

GOPATH=$(shell pwd)/_workspace/
SRC_PATH=$(GOPATH)src/$(PROJECT)
GOBUILD=GOPATH=$(GOPATH) go build
GOCLEAN=GOPATH=$(GOPATH) go clean

BUILD_DIR=build
BUILD_BIN_DIR=$(BUILD_DIR)/bin
BIN_PATH=$(DESTDIR)/bin/
BIN_FILE_LIST=git-comment git-comment-grep git-comment-log git-comment-remote git-comment-web
BIN_BUILD_CMD=$(GOBUILD) -ldflags "-X main.buildVersion=$(VERSION)"

MAN_PATH=$(DESTDIR)/man/man1/
MAN_BUILD_DIR=$(BUILD_DIR)/man/
MAN_TITLE=Git Comment Manual
MAN_CMD=pod2man --center="$(MAN_TITLE)" --release="$(VERSION)"

all: bootstrap build install

default: build

bootstrap: env
	GOPATH=$(GOPATH) go get \
				 gopkg.in/libgit2/git2go.v23 \
				 github.com/stvp/assert \
	       github.com/cevaris/ordered_map \
	       gopkg.in/alecthomas/kingpin.v2 \
	       github.com/kylef/result.go/src/result \
	       github.com/blevesearch/bleve \
	       github.com/blang/semver

bootstrap_osx:
	brew install libgit2

build: copy
	$(GOBUILD) $(PROJECT)
	mkdir -p $(BUILD_BIN_DIR)
	$(foreach bin,$(BIN_FILE_LIST),$(BIN_BUILD_CMD) -o $(BUILD_BIN_DIR)/$(bin) bin/$(bin).go;)

ci: bootstrap test

clean: env
	$(GOCLEAN) -i -x $(PROJECT) || true
	rm -rf $(SRC_PATH) $(BUILD_DIR)

copy: env
	$(foreach pack,$(PACKAGES),install -d $(SRC_PATH)/$(pack);)
	install $(PROJECT)/*.go $(SRC_PATH)
	$(foreach pack,$(PACKAGES),install $(PROJECT)/$(pack)/*.go $(SRC_PATH)/$(pack);)

dep: env
	rm -r $(GOPATH)*/*/*/.git || true
	git add $(GOPATH)

deploy_website:
	git checkout -B gh-pages
	git filter-branch -f --subdirectory-filter docs/git-comment.com
	git clean -df
	# git push -f origin gh-pages

doc:
	mkdir -p $(MAN_BUILD_DIR)
	$(foreach bin,$(BIN_FILE_LIST), $(MAN_CMD) docs/man/$(bin).pod > $(MAN_BUILD_DIR)$(bin).1;)

env:
	install -d $(GOPATH)

install: doc
	$(foreach bin,$(BIN_FILE_LIST), \
		chown root:admin $(MAN_BUILD_DIR)$(bin).1; \
		chmod 444 $(MAN_BUILD_DIR)$(bin).1; \
		install $(BUILD_BIN_DIR)/$(bin) $(BIN_PATH)$(bin); \
		gzip -f $(MAN_BUILD_DIR)$(bin).1; \
		install -C $(MAN_BUILD_DIR)$(bin).1.gz $(MAN_PATH)$(bin).1.gz;)
	rm -r $(MAN_BUILD_DIR)

uninstall:
	$(foreach bin,$(BIN_FILE_LIST), rm $(MAN_PATH)$(bin).1.gz $(BIN_PATH)$(bin);)

test: copy
	go test $(PROJECT) $(foreach pkg,$(PACKAGES),$(PROJECT)/$(pkg));
