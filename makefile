# Copyright (c) 2018 Anton Semjonov
# Licensed under the MIT License

OUTPUT := aenker

# compile statically linked binary
.PHONY: build
build : $(OUTPUT)
$(OUTPUT) : $(shell find * -type f -name '*.go') go.mod go.sum
	CGO_ENABLED=0 go build -ldflags '-s -w' -o "$(OUTPUT)"

# ansemjo/makerelease targets
include makerelease

# make a release / cross-compile with mkr
release:
	git archive --prefix=./ HEAD | mkr rl



PREFIX := $(shell [ $$(id -u) -eq 0 ] && echo /usr/local || echo ~/.local)
INSTALLED := $(PREFIX)/bin/$(OUTPUT)
MANUALS := $(PREFIX)/share/man

# install OUTPUT and docs
install : $(INSTALLED) $(MANUALS)/man1/$(OUTPUT).1
$(INSTALLED) : $(OUTPUT)
	install -m 755 $< $@

# generate local docs
docs : docs/$(OUTPUT).md
docs/$(OUTPUT).md : $(OUTPUT)
	mkdir -p docs
	./$< gen manual -d docs man
	./$< gen manual -d docs markdown

# generate manuals
manuals : $(MANUALS)/man1/$(OUTPUT).1
$(MANUALS)/man1/$(OUTPUT).1 : $(OUTPUT)
	./$< gen manual -d $(MANUALS)
	@echo "# add this to your ~/.bashrc:"
	@echo ". <($< gen completion)"
	@echo "# or add global bash completions:"
	@echo "$< gen completion > /usr/share/bash-completion/completions/$<"

# clean untracked files and directories
clean :
	git clean -fdx
