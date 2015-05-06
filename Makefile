PROJECT=git-comment
BIN_PATH=/usr/local/bin/$(PROJECT)
MAN_PATH=/usr/share/man/man1/$(PROJECT).1
SRC_PATH=$(GOPATH)/src/$(PROJECT)
MAN_TMP_PATH=$(PROJECT).1
MAN_ZIP_PATH=$(PROJECT).1.gz

default: build

bootstrap:
	mkdir -p man src/$(PROJECT)
	touch README.md LICENSE CONTRIBUTING.md CHANGELOG.md TODO.md Makefile man/$(PROJECT).1

build: copy
	go build $(PROJECT)

clean:
	go clean -i -x $(PROJECT)
	rm -rf $(SRC_PATH)
	rm $(PROJECT)

copy:
	install -d $(SRC_PATH)
	install src/$(PROJECT)/* $(SRC_PATH)

install:
	install $(PROJECT) $(BIN_PATH)
	cp man/$(PROJECT).1 $(MAN_TMP_PATH)
	chown root:admin $(MAN_TMP_PATH)
	chmod 444 $(MAN_TMP_PATH)
	tar -czf $(MAN_ZIP_PATH) $(MAN_TMP_PATH)
	install -C $(MAN_TMP_PATH) $(MAN_PATH)
	rm $(MAN_TMP_PATH) $(MAN_ZIP_PATH)

uninstall:
	rm $(MAN_PATH) $(BIN_PATH)

test: copy
	go test $(PROJECT)
