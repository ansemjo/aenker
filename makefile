# Copyright (c) 2018 Anton Semjonov
# Licensed under the MIT License

# ---------- build ----------

NAME    := aenker
VERSION := $(shell sh version.sh describe)

# env and flags to build static binary with embedded version
GO_BUILD_ENV   := CGO_ENABLED=0
GO_BUILD_FLAGS := -ldflags='-s -w -X main.Version=$(VERSION)'

# build binary for host system
.PHONY: build
build : $(NAME)
$(NAME) : $(shell find * -type f -name '*.go') go.mod go.sum
	env $(GO_BUILD_ENV) go build $(GO_BUILD_FLAGS) -o $@

# cross-compile binaries with gox
.PHONY: release
release :
	env $(GO_BUILD_ENV) gox $(GO_BUILD_FLAGS) -output='$@/$(NAME)-{{.OS}}-{{.Arch}}'

# run tests verbosely
.PHONY: test
test :
	go test ./... -timeout 30s -cover -v

# ---------- install ----------

# installation directories
DESTDIR        :=
PREFIX         := /usr
BINARY_DIR     := $(DESTDIR)$(PREFIX)/bin
MANUAL_DIR     := $(DESTDIR)$(PREFIX)/share/man
COMPLETION_DIR := $(DESTDIR)$(PREFIX)/share/bash-completion/completions

# install binary and manuals
.PHONY: install
install : \
	$(BINARY_DIR)/$(NAME) \
	$(MANUAL_DIR)/man1/$(NAME).1 \
	$(COMPLETION_DIR)/$(NAME)

$(BINARY_DIR)/$(NAME) : $(NAME)
	install -m 755 -D $< $@

$(MANUAL_DIR)/man1/$(NAME).1 : $(NAME)
	install -m 755 -d $(MANUAL_DIR)
	./$< docs manual man -d $(MANUAL_DIR)

$(COMPLETION_DIR)/$(NAME) : $(NAME)
	install -m 755 -d $(COMPLETION_DIR)
	./$< docs completion bash -o $@

# ---------- packaging ----------

# package metadata
PKGNAME     := $(NAME)
PKGVERSION  := $(shell echo $(VERSION) | sed s/-/./ )
PKGAUTHOR   := 'ansemjo <anton@semjonov.de'
PKGLICENSE  := MIT
PKGURL      := https://github.com/ansemjo/$(PKGNAME)
PKGFORMATS  := rpm deb apk
PKGARCH     := $(shell uname -m)

# how to execute fpm
FPM = podman run --rm --net none -v $$PWD:/src -w /src ansemjo/fpm:alpine

# build a package
.PHONY: package-%
package-% :
	make --no-print-directory install DESTDIR=package
	mkdir -p release
	$(FPM) -s dir -t $* -f --chdir package \
		--name $(PKGNAME) \
		--version $(PKGVERSION) \
		--maintainer $(PKGAUTHOR) \
		--license $(PKGLICENSE) \
		--url $(PKGURL) \
		--architecture $(PKGARCH) \
		--package release/$(PKGNAME)-$(PKGVERSION)-$(PKGARCH).$*

# build all package formats with fpm
.PHONY: packages
packages : $(addprefix package-,$(PKGFORMATS))

# ---------- misc ----------

# generate local docs
docs : $(NAME)
	mkdir -p $@
	./$< docs manual -d $@ markdown

# clean untracked files and directories
.PHONY: clean
clean :
	git clean -fdx
