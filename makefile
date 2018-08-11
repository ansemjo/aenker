# mini makefile to build binary with build.go
# and install it in ~/.local/bin
#
# this compiles aenker as a static binary and
# adds a version tag from git. installing via
# 'go install' works aswell but misses the
# above two features

.PHONY  : default install build clean docs uninstall

BINARY := aenker
PREFIX := ~/.local
INSTALLED := $(PREFIX)/bin/$(BINARY)
MANUALS := $(PREFIX)/share/man

default : build

# clean untracked files and directories
clean :
	git clean -dfx

# install vendored packages with https://github.com/golang/dep
vendor :
	dep ensure

# compile static binary
build : $(BINARY)
$(BINARY) : vendor $(shell find * -type f -name '*.go')
	go run build.go -o $@
	command -V upx >/dev/null && upx $@

# install binary and docs
install : $(INSTALLED)
$(INSTALLED) : $(BINARY) $(MANUALS)/man1/$(BINARY).1
	install -m 755 $< $@

# generate documentation
docs : $(MANUALS)/man1/$(BINARY).1
$(MANUALS)/man1/$(BINARY).1 : $(BINARY)
	mkdir -p docs
	./$< gen manual -d docs
	./$< gen manual -d docs markdown
	./$< gen manual -d $(MANUALS)
	@echo "# add this to your ~/.bashrc:"
	@echo ". <(aenker gen completion)"

# attempt to remove installed files
uninstall :
	rm -fv $(INSTALLED)
	rm -fv $(MANUALS)/man1/$(BINARY)*.1