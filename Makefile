# Target directory for installing git-comment
DESTDIR := /usr/local
# Subdirectory within $(DESTDIR) for installing the binaries
DESTBIN := $(DESTDIR)/bin
# Subdirectory within $(DESTDIR) for installing the manual pages
DESTMAN := $(DESTDIR)/share/man/man1

# Command to install a file into a path without changing the modification
# date if not needed. Default is BSD-style arguments.
INSTALLCMD    := install -C
# Command to install a directory to a path. Default is BSD-style arguments.
INSTALLDIRCMD := install -d

PROJECT=libgitcomment
# Subpackages within the libgitcomment library
PACKAGES=exec log git search
# Current version of the git-comment tool
VERSION=$(shell cat VERSION)
# Go libraries on which libgitcomment and the git-comment tool depend
DEPENDENCIES=gopkg.in/libgit2/git2go.v23 \
	github.com/stvp/assert \
	github.com/cevaris/ordered_map \
  gopkg.in/alecthomas/kingpin.v2 \
  github.com/kylef/result.go/src/result \
  github.com/blevesearch/bleve \
  github.com/blang/semver
# List of binary file names within the git-comment suite
BIN_FILES=$(basename $(shell ls bin))
# List of non-test source files within libgitcomment
SRC_FILES=$(filter-out test,$(shell git ls-files "$(PROJECT)/**.go"))
# List of test files within libgitcomment
TEST_FILES=$(shell git ls-files "$(PROJECT)/**_test.go")
# Directory containing manual page source files
MANSRC=docs/man
# List of source files in the manual directory
MAN_FILES=$(foreach bin,$(BIN_FILES),$(bin).pod)

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
# Vendored sources and temporary build directory
GOPATH=$(shell pwd)/_workspace
# Command to execute go with the default arguments for the project
GO=GOPATH="$(GOPATH)" go
# Temporary source directory for libgitcomment within the temporary build
# directory
GOPATHSRC=$(GOPATH)/src/$(PROJECT)
# List of files within libgitcomment copied within the temporary build
# directory
GOPATHSRC_FILES=$(addprefix $(GOPATH)/src/,$(SRC_FILES))
# List of test files within libgitcomment copied within the temporary build
# directory
GOPATHSRC_TESTS=$(addprefix $(GOPATH)/src/,$(TEST_FILES))
# Target directory for compiled dependent libraries
GOPATHPKG=$(GOPATH)/pkg/$(GOOS)_$(GOARCH)
# List of all dependent libraries
GOPATHPKG_DEPS=$(foreach dep,$(DEPENDENCIES),$(GOPATHPKG)/$(dep).a)

# Output directory for project artifacts
BUILD_DIR=build
# Output directory for compiled binaries
BUILD_BIN_DIR=$(BUILD_DIR)/bin
# Output directory for manual pages
BUILD_MAN_DIR=$(BUILD_DIR)/man
# List of compiled binaries
BUILD_BIN_FILES=$(foreach bin,$(BIN_FILES),$(BUILD_BIN_DIR)/$(bin))
# List of manual pages
BUILD_MAN_FILES=$(foreach bin,$(BIN_FILES),$(BUILD_MAN_DIR)/$(bin).1)

# Title of the git-comment manual
MAN_TITLE=Git Comment Manual
# Command to build the manual pages denoting the title and release number
MAN_CMD=pod2man --center="$(MAN_TITLE)" --release="$(VERSION)"

default: build

# Build the core library 'libgitcomment' as well as the tool binaries
# listed in bin/.
build: $(BUILD_BIN_FILES)

# Install necessary dependencies for building git-comment on Ubuntu 12.04
ci_deps: apt_deps src_libgit2

# Install necessary dependencies via apt-get for building from an empty env
# on Ubuntu 12.04
apt_deps:
	apt-get update
	apt-get install -y cmake pkg-config

# Download and install libgit2 0.23.2 from source. Depends on cmake, pkg-config.
src_libgit2:
	wget -O libgit2.tar.gz https://github.com/libgit2/libgit2/archive/v0.23.2.tar.gz
	tar xzvf libgit2.tar.gz
	cd libgit2-0.23.2 && cmake . && make && make install
	ldconfig

# $(BUILD_BIN_DIR) is the target directory for the compiled binaries of the
# tools listed in the bin/ directory. Building each binary depends on the
# go library dependencies being built and having the latest versions of the
# files in bin/.
# This target first ensures the build directory exists or creates it, then
# builds each binary with a flag specifying the version of the git-comment
# project.
$(BUILD_BIN_DIR)/%: $(GOPATHSRC_FILES) $(GOPATHPKG_DEPS) bin/%.go
	$(INSTALLDIRCMD) $(BUILD_BIN_DIR)
	$(GO) build -ldflags "-X main.buildVersion=$(VERSION)" -o "$(BUILD_BIN_DIR)/$*" bin/$*.go

# $(GOPATHSRC) is a temporary build directory within the local $GOPATH for
# the source files in libgitcomment/. Building the library depends on the
# go library dependencies being built and having the latest versions of the
# files in libgitcomment/.
# This target ensures the temporary build directory exists or creates it, then
# installs changed source files into it.
$(GOPATHSRC)/%.go: $(GOPATHPKG_DEPS) $(PROJECT)/%.go
	$(INSTALLDIRCMD) "$(GOPATHSRC)/$(dir $*)"
	$(INSTALLCMD) "$(PROJECT)/$*.go" "$(GOPATHSRC)/$*.go"

# $(GOPATHPKG) is the compiled binary path within the local $GOPATH.
# This target builds a library file for any repository specified by the
# DEPENDENCIES list, then removes git metadata so the downloaded source files
# can be checked into source control.
$(GOPATHPKG)/%.a:
	$(GO) get $*
	rm -rf "$(GOPATH)/src/$*/.git"

# Remove compiled files and build directories for the project so it can be
# rebuilt from a clean slate.
clean:
	$(GO) clean $(PROJECT) || true
	rm -rf "$(GOPATHSRC)" "$(BUILD_DIR)"

# Convert the POD-format manual page source files into *roff output. Depends
# on the documentation source files in $(MANSRC)
$(BUILD_MAN_DIR)/%.1: $(MANSRC)/%.pod
	@$(INSTALLDIRCMD) $(BUILD_MAN_DIR)
	$(MAN_CMD) $(MANSRC)/$*.pod > $(BUILD_MAN_DIR)/$*.1
	chmod 444 $(BUILD_MAN_DIR)/$*.1

# Generate *roff-format manual pages for each of the tool binaries manual
# sources in $(MANSRC)
doc: $(BUILD_MAN_FILES)

# Create the target directory for installing tool binaries if it does not
# exist
$(DESTBIN):
	@$(INSTALLDIRCMD) $(DESTBIN)

# Create the target directory for installing tool manuals if it does not
# exist
$(DESTMAN):
	@$(INSTALLDIRCMD) $(DESTMAN)

# Install git-comment into the preferred path on the host machine. Depends
# on building the tool binaries and manual
install: $(DESTBIN) $(DESTMAN) $(BUILD_BIN_FILES) $(BUILD_MAN_FILES)
	@$(INSTALLCMD) $(BUILD_BIN_FILES) $(DESTBIN)
	@$(INSTALLCMD) $(BUILD_MAN_FILES) $(DESTMAN)
	@echo Successfully installed git-comment.

# Remove git-comment from the preferred path on the host machine
uninstall:
	rm $(foreach bin,$(BIN_FILES), $(DESTMAN)/$(bin).1 $(DESTBIN)/$(bin));

# Run the unit test suite on the libgitcomment library
test: $(GOPATHSRC_FILES) $(GOPATHSRC_TESTS)
	$(GO) test $(PROJECT) $(foreach pkg,$(PACKAGES),$(PROJECT)/$(pkg));
