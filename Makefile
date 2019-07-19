PKGNAME          = otterchat
DESTDIR         ?=
PREFIX          ?= /usr
VERSION         ?=

BINDIR           = $(PREFIX)/bin
CONFIGDIR        = $(HOME)/.config
PLUGINDIR        = $(CONFIGDIR)/$(PKGNAME)/plugins
GOBIN            = $(GOPATH)/bin

GOCC             = go
GOLDFLAGS        = -ldflags "-s -w"
GOTAGS           = --tags "chat irc"

GOBUILD          = $(GOCC) build -v $(GOLDFLAGS) $(GOTAGS)
GOPLUGIN         = $(GOCC) build -buildmode=plugin
GOTEST           = $(GOCC) test
GOMOD            = $(GOCC) mod
GOLINT           = $(GOBIN)/golint -set_exit_status
GOVET            = $(GOCC) vet
GOFMT            = $(GOCC) fmt -x

include Makefile.waterlog

all: build plugins

build:
	@$(call stage,BUILD)
	@$(call task,Building executable...)
	@$(GOBUILD)
	@$(call task,Creating Plugin Directory...)
	@mkdir -p $(PLUGINDIR)
	@$(call task,Building IRC plugin...)
	@$(GOPLUGIN) -o irc.so plugins/irc.go
	@install -D -m 00755 irc.so $(PLUGINDIR)/irc.so
	@$(call pass,BUILD)

test:
	@$(call stage,TEST)
	@$(GOTEST) -cover ./...
	@$(call pass,TEST)

validate:
	@$(call stage,FORMAT)
	@$(GOFMT) ./...
	@$(call pass,FORMAT)
	@$(call stage,VET)
	@$(call task,Running 'go vet'...)
	@$(GOVET) ./...
	@$(call pass,VET)
	@$(call stage,LINT)
	@$(call task,Running 'golint'...)
	@$(GOLINT) `go list ./... | grep -v vendor`
	@$(call pass,LINT)

tidy:
	@$(call stage,TIDY)
	@$(GOMOD) tidy
	@$(call pass,TIDY)

vendor:
	@$(call stage,VENDOR)
	@$(GOMOD) vendor
	@$(call pass,VENDOR)

install:
	@$(call stage,INSTALL)
	install -D -m 00755 $(PKGNAME) $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,INSTALL)

uninstall:
	@$(call stage,UNINSTALL)
	rm -f $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,UNINSTALL)

clean:
	@$(call stage,CLEAN)
	@$(call task,Removing executables...)
	@rm -f $(PKGNAME)
	@rm -f irc.so
	@rm -rf vendor
	@rm -f *.tar.gz
	@$(call pass,CLEAN)

package: clean tidy vendor
	@$(call stage,PACKAGE)
ifndef VERSION
	@echo "VERSION is not defined"
	@$(call fail,PACKAGE)
else
	@$(call task,Building archive...)
	@tar --exclude='.git' --exclude='*.tar.gz' -zcvf $(PKGNAME)-v$(VERSION).tar.gz ../$(PKGNAME)
	@$(call task,tagging version...)
	#@git tag -a v$(VERSION) -m "$(PKGNAME) $(VERSION)"
	@$(call pass,PACKAGE)
endif