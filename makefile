# Copyright (c) 2018 Anton Semjonov
# Licensed under the MIT License

# mini makefile to build binary with build.go
# and install it in ~/.local/bin
#
# this compiles aenker as a static binary and
# adds a version tag from git. installing via
# 'go install' works aswell but misses the
# above two features

.PHONY  : default install build clean manuals uninstall docs

BINARY := aenker
PREFIX := $(shell [ $$(id -u) -eq 0 ] && echo /usr/local || echo ~/.local)
INSTALLED := $(PREFIX)/bin/$(BINARY)
MANUALS := $(PREFIX)/share/man
GOFILES := $(shell find * -type f -name '*.go')
TMPGOPATH := /tmp/aenker-build-tmpgopath

# install vendored packages with https://github.com/golang/dep
# compile static binary
build : $(BINARY)
$(BINARY) : $(GOFILES)
	vgo mod vendor
	go run build.go -o $@ --tempdir $(TMPGOPATH)
	./$@ --version
	sha256sum --tag $@

# compress binary with upx
compress : $(BINARY)
	upx $<

# prepare for ansemjo/makerelease
mkrelease-prepare:
	go mod vendor

mkrelease-targets:
	@bash -c 'echo {linux,darwin}/{386,amd64} linux/arm{,64} {free,open}bsd/{386,amd64,arm}'

mkrelease:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) \
		go build -o $(RELEASEDIR)/$(BINARY)-$(OS)-$(ARCH)

mkrelease-finish:
	upx $(RELEASEDIR)/* || true
	cd $(RELEASEDIR) && sha256sum $(BINARY)-*-* | tee sha256sums

# install binary and docs
install : $(INSTALLED) $(MANUALS)/man1/$(BINARY).1
$(INSTALLED) : $(BINARY)
	install -m 755 $< $@

# generate local docs
docs : docs/$(BINARY).md
docs/$(BINARY).md : $(BINARY)
	mkdir -p docs
	./$< gen manual -d docs man
	./$< gen manual -d docs markdown

# generate manuals
manuals : $(MANUALS)/man1/$(BINARY).1
$(MANUALS)/man1/$(BINARY).1 : $(BINARY)
	./$< gen manual -d $(MANUALS)
	@echo "# add this to your ~/.bashrc:"
	@echo ". <($< gen completion)"
	@echo "# or add global bash completions:"
	@echo "$< gen completion > /usr/share/bash-completion/completions/$<"

# clean untracked files and directories
clean :
	git clean -dfx

# attempt to remove installed files
uninstall :
	rm -fv $(INSTALLED)
	rm -fv $(MANUALS)/man1/$(BINARY)*.1
