PROJECT=git_comment
PACKAGES=exec log git search
VERSION=$(shell cat VERSION)
SRC_PATH=$(GOPATH)src/$(PROJECT)
BUILD_DIR=build/bin/
BIN_PATH=/usr/local/bin/
BIN_FILE_LIST=git-comment git-comment-grep git-comment-log git-comment-web
BIN_BUILD_CMD=go build -ldflags "-X main.buildVersion $(VERSION)"
MAN_PATH=/usr/local/man/man1/
MAN_TMP_PATH=build/doc/
MAN_TITLE=Git Comment Manual
MAN_CMD=pod2man --center="$(MAN_TITLE)" --release="$(VERSION)"

default: build

bootstrap:
	brew install libgit2
	go get github.com/libgit2/git2go
	go get github.com/stvp/assert
	go get github.com/cevaris/ordered_map
	go get gopkg.in/alecthomas/kingpin.v2
	go get github.com/kylef/result.go/src/result
	go get github.com/blevesearch/bleve
	go get github.com/blang/semver

build: copy
	go build $(PROJECT)
	mkdir -p $(BUILD_DIR)
	$(foreach bin,$(BIN_FILE_LIST),$(BIN_BUILD_CMD) -o $(BUILD_DIR)/$(bin) bin/$(bin).go;)

clean:
	go clean -i -x $(PROJECT)
	rm -rf $(SRC_PATH) $(BUILD_DIR)

copy:
	$(foreach pack,$(PACKAGES),install -d $(SRC_PATH)/$(pack);)
	install $(PROJECT)/*.go $(SRC_PATH)
	$(foreach pack,$(PACKAGES),install $(PROJECT)/$(pack)/*.go $(SRC_PATH)/$(pack);)

doc:
	$(foreach bin,$(BIN_FILE_LIST), $(MAN_CMD) man/$(bin).pod > man/$(bin).1;)

install: doc
	mkdir -p $(MAN_TMP_PATH)
	$(foreach bin,$(BIN_FILE_LIST), \
		cp man/$(bin).1 $(MAN_TMP_PATH)$(bin).1; \
		chown root:admin $(MAN_TMP_PATH)$(bin).1; \
		chmod 444 $(MAN_TMP_PATH)$(bin).1; \
		install $(BUILD_DIR)/$(bin) $(BIN_PATH)$(bin); \
		gzip -f $(MAN_TMP_PATH)$(bin).1; \
		install -C $(MAN_TMP_PATH)$(bin).1.gz $(MAN_PATH)$(bin).1.gz;)
	rm -r $(MAN_TMP_PATH)

uninstall:
	$(foreach bin,$(BIN_FILE_LIST), rm $(MAN_PATH)$(bin).1 $(BIN_PATH)$(bin);)

test: copy
	go test $(PROJECT)
	$(foreach pkg,$(PACKAGES),go test $(PROJECT)/$(pkg);)
