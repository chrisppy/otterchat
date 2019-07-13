PKGNAME = otterchat
DESTDIR ?=
PREFIX  ?= /usr
XDG_CONFIG_HOME = $(HOME)/.config

all: build

build:
	go build
	mkdir -p $(XDG_CONFIG_HOME)/$(PKGNAME)/plugins
	go build -buildmode=plugin -o irc.so plugins/irc.go
	install -D -m 00755 irc.so $(XDG_CONFIG_HOME)/$(PKGNAME)/plugins/irc.so
test:
	go test -v ./...

cover:
	go test -cover ./...

tidy:
	go mod tidy

vendor:
	go mod vendor

install:
	install -D -m 00755 $(PKGNAME) $(DESTDIR)$(PREFIX)/bin/$(PKGNAME)

uninstall:
	rm $(DESTDIR)$(PREFIX)/bin/$(PKGNAME)

clean:
	rm -rf $(XDG_CONFIG_HOME)/$(PKGNAME)
	rm $(PKGNAME)