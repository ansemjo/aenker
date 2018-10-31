# Copyright (c) 2018 Anton Semjonov
# Licensed under the MIT License

# ---------- build ----------

# compile statically linked binary
.PHONY: build
OUTPUT := aenker
build : $(OUTPUT)
$(OUTPUT) : $(shell find * -type f -name '*.go') go.mod go.sum
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags '-s -w' -o $(OUTPUT)

# makerelease targets for reproducible builds, ansemjo/makerelease
.PHONY : mkrelease-prepare mkrelease mkrelease-finish
mkrelease-prepare :
	go mod download

EXT := $(if $(findstring windows,$(OS)),.exe)
mkrelease :
	OUTPUT=$(RELEASEDIR)/$(OUTPUT)-$(OS)-$(ARCH)$(EXT) make --no-print-directory build

mkrelease-finish :
	upx -q $$(find $(RELEASEDIR)/* ! -name '*bsd-a*')
	printf "# built with %s in %s\n" "$$MKR_VERSION" "$$MKR_IMAGE" > $(RELEASEDIR)/SHA256SUMS
	cd $(RELEASEDIR) && sha256sum $(OUTPUT)-*-* | tee -a SHA256SUMS

# make a release / cross-compile with mkr
release:
	git archive --prefix=./ HEAD | mkr rl $(MKRARGS)

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
	$(BINARY_DIR)/$(OUTPUT) \
	$(MANUAL_DIR)/man1/$(OUTPUT).1 \
	$(COMPLETION_DIR)/$(OUTPUT)

$(BINARY_DIR)/$(OUTPUT) : $(OUTPUT)
	install -m 755 -D $< $@

$(MANUAL_DIR)/man1/$(OUTPUT).1 : $(OUTPUT)
	install -m 755 -d $(MANUAL_DIR)
	./$< docs manual man -d $(MANUAL_DIR)

$(COMPLETION_DIR)/$(OUTPUT) : $(OUTPUT)
	install -m 755 -d $(COMPLETION_DIR)
	./$< docs completion bash -o $@

# ---------- packaging ----------

# package metadata
PKGNAME     := aenker
PKGVERSION  := $(shell sh version.sh describe | sed s/-/./ )
PKGAUTHOR   := 'ansemjo <anton@semjonov.de'
PKGLICENSE  := MIT
PKGURL      := https://github.com/ansemjo/$(PKGNAME)
PKGFORMATS  := rpm deb apk pacman
PACKAGEDIR  := package

# how to execute fpm
FPM = podman run --rm --net none -v $$PWD:/build -w /build ansemjo/fpm:alpine

# build a package
.PHONY: package-%
package-% :
	make --no-print-directory install DESTDIR=$(PACKAGEDIR)
	$(FPM) -s dir -t $* -f --chdir $(PACKAGEDIR) \
		--name $(PKGNAME) \
		--version $(PKGVERSION) \
		--maintainer $(PKGAUTHOR) \
		--license $(PKGLICENSE) \
		--url $(PKGURL)

# build all package formats with fpm
.PHONY: packages
packages : $(addprefix package-,$(PKGFORMATS))

# ---------- misc ----------

# generate local docs
docs : $(OUTPUT)
	mkdir -p $@
	./$< docs manual -d $@ markdown

# clean untracked files and directories
.PHONY: clean
clean :
	git clean -fdx
