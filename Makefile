DESTDIR := /usr/local
DESTBIN=$(DESTDIR)/bin/

PROJECT=libgitcomment
PACKAGES=exec log git search
VERSION=$(shell cat VERSION)
BIN_FILES=git-comment git-comment-grep git-comment-log git-comment-remote git-comment-web
SRC_FILES=$(PROJECT)/*.go $(foreach pkg,$(PACKAGES), $(PROJECT)/$(pkg)/*.go)

GOPATH=$(shell pwd)/_workspace/
SRC_PATH=$(GOPATH)src/$(PROJECT)
GOBUILD=GOPATH=$(GOPATH) go build
GOCLEAN=GOPATH=$(GOPATH) go clean
BIN_BUILD_CMD=$(GOBUILD) -ldflags "-X main.buildVersion=$(VERSION)"

BUILD_DIR=build
BUILD_BIN_DIR=$(BUILD_DIR)/bin
BUILD_BIN_FILES=$(foreach bin,$(BIN_FILES),$(BUILD_BIN_DIR)/$(bin))

MAN_PATH=$(DESTDIR)/man/man1/
MAN_BUILD_DIR=$(BUILD_DIR)/man/
MAN_TITLE=Git Comment Manual
MAN_CMD=pod2man --center="$(MAN_TITLE)" --release="$(VERSION)"

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

gobuild: copy
	$(GOBUILD) $(PROJECT)

build: gobuild $(BUILD_BIN_FILES)

build/bin/%: bin/%.go
	@mkdir -p $(BUILD_BIN_DIR)
	$(BIN_BUILD_CMD) -o $(BUILD_BIN_DIR)/$* bin/$*.go

$(SRC_PATH)/%/*.go: $(PROJECT)/%/*.go
	install -d $(PROJECT)/$* $(SRC_PATH)/$*
	install $(PROJECT)/$*/*.go $(SRC_PATH)/$*

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
	$(foreach bin,$(BIN_FILES), $(MAN_CMD) docs/man/$(bin).pod > $(MAN_BUILD_DIR)$(bin).1;)

env:
	install -d $(GOPATH)

install: bootstrap build doc
	$(foreach bin,$(BIN_FILES), \
		chown root:admin $(MAN_BUILD_DIR)$(bin).1; \
		chmod 444 $(MAN_BUILD_DIR)$(bin).1; \
		install $(BUILD_BIN_DIR)/$(bin) $(DESTBIN)$(bin); \
		gzip -f $(MAN_BUILD_DIR)$(bin).1; \
		install -C $(MAN_BUILD_DIR)$(bin).1.gz $(MAN_PATH)$(bin).1.gz;)
	rm -r $(MAN_BUILD_DIR)

uninstall:
	$(foreach bin,$(BIN_FILES), rm $(MAN_PATH)$(bin).1.gz $(DESTBIN)$(bin);)

test: copy
	go test $(PROJECT) $(foreach pkg,$(PACKAGES),$(PROJECT)/$(pkg));
