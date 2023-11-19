PKGNAME = pfetch
ARCHIVE = https://github.com/dylanaraps/pfetch/archive/refs/tags/0.6.0.tar.gz
WORKDIR = /tmp/$(PKGNAME)-build
DOWNLOADER = /usr/bin/wget -O $(WORKDIR)/$(PKGNAME)-archive.tar.gz
DESTDIR = /usr/bin
VERSION = 0.6.0
DEPS = bash

all: depinfo dl extract install clean

dl:
	mkdir -pv $(WORKDIR)
	$(DOWNLOADER) $(ARCHIVE)

extract:
	tar -xvf $(WORKDIR)/$(PKGNAME)-archive.tar.gz -C $(WORKDIR) --strip-components=1

install:
	mkdir -pv $(DESTDIR)
	install -v -o root -g root -m 0755 -p $(WORKDIR)/pfetch $(DESTDIR)/$(PKGNAME)

uninstall:
	rm -rvf $(DESTDIR)/$(PKGNAME)

version:
	echo $(PKGNAME): Version $(VERSION)

depinfo:
	@echo -n "The package $(PKGNAME) marks [$(DEPS)] as dependencies, rpkgm does not support dependency resolution and installation yet, you will need to install them manually.\n"

clean:
	rm -rvf $(WORKDIR)