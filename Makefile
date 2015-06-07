PROJECT=git_comment
VERSION=$(shell cat VERSION)
SRC_PATH=$(GOPATH)/src/$(PROJECT)
BIN_PATH=/usr/local/bin/
BIN_FILE_LIST=git-comment git-comment-grep git-comment-log
BIN_BUILD_CMD=go build -ldflags "-X main.buildVersion $(VERSION)"
MAN_PATH=/usr/share/man/man1/
MAN_TMP_PATH=build/doc/
MAN_ZIP_PATH=$(PROJECT).1.gz
MAN_TITLE=Git Comment Manual
MAN_CMD=pod2man --center="$(MAN_TITLE)" --release="$(VERSION)"

default: build

bootstrap:
	brew install libgit2
	go get gopkg.in/libgit2/git2go.v22
	go get github.com/stvp/assert
	go get github.com/cevaris/ordered_map
	go get github.com/droundy/goopt

build: copy
	go build $(PROJECT)
	$(foreach bin,$(BIN_FILE_LIST),$(BIN_BUILD_CMD) bin/$(bin).go;)

clean:
	go clean -i -x $(PROJECT)
	rm -rf $(SRC_PATH)
	$(foreach bin,$(BIN_FILE_LIST),rm $(bin);)

copy:
	install -d $(SRC_PATH)
	install src/$(PROJECT)/* $(SRC_PATH)

doc:
	$(foreach bin,$(BIN_FILE_LIST), $(MAN_CMD) man/$(bin).pod > man/$(bin).1;)

install: doc
	mkdir -p $(MAN_TMP_PATH)
	$(foreach bin,$(BIN_FILE_LIST), \
		cp man/$(bin).1 $(MAN_TMP_PATH)$(bin).1; \
		chown root:admin $(MAN_TMP_PATH)$(bin).1; \
		chmod 444 $(MAN_TMP_PATH)$(bin).1; \
		install $(bin) $(BIN_PATH)$(bin); \
		gzip -f $(MAN_TMP_PATH)$(bin).1; \
		install -C $(MAN_TMP_PATH)$(bin).1.gz $(MAN_PATH)$(bin).1.gz;)
	rm -r $(MAN_TMP_PATH)

uninstall:
	$(foreach bin,$(BIN_FILE_LIST), rm $(MAN_PATH)$(bin).1 $(BIN_PATH)$(bin);)

test: copy
	go test $(PROJECT)
