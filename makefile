# mini makefile to build binary with build.go
# and install it in ~/.local/bin
#
# this compiles aenker as a static binary and
# adds a version tag from git. installing via
# 'go install' works aswell but misses the
# above two features

BINARY := aenker
PREFIX := ~/.local
INSTALLED := $(PREFIX)/bin/$(BINARY)

.PHONY  : install build
build   : $(BINARY)
install : $(INSTALLED)

# install vendored packages with https://github.com/golang/dep
vendor :
	dep ensure

# compile binary
$(BINARY) : vendor
	go run build.go -o $@

# install binary
$(INSTALLED) : $(BINARY)
	install -m 755 $< $@