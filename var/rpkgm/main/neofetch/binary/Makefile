PKGNAME = neofetch
ARCHIVE = https://github.com/dylanaraps/neofetch/archive/refs/tags/7.1.0.tar.gz
WORKDIR = /tmp/$(PKGNAME)-build
DOWNLOADER = /usr/bin/wget -O $(WORKDIR)/$(PKGNAME)-archive.tar.gz
DESTDIR = /usr/bin
MANDIR = /usr/share/man/man1
VERSION = 7.1.0
DEPS = bash

all: depinfo dl extract install clean

dl:
	mkdir -pv $(WORKDIR)
	$(DOWNLOADER) $(ARCHIVE)

extract:
	tar -xvf $(WORKDIR)/$(PKGNAME)-archive.tar.gz -C $(WORKDIR) --strip-components=1

install:
	mkdir -pv $(DESTDIR)
	mkdir -pv $(MANDIR)
	install -v -o root -g root -m 0755 -p $(WORKDIR)/neofetch $(DESTDIR)/$(PKGNAME)
	install -v -o root -g root -m 0644 -p $(WORKDIR)/neofetch.1 $(MANDIR)

uninstall:
	rm -rvf $(DESTDIR)/$(PKGNAME)
	rm -rvf $(MANDIR)/neofetch.1*

version:
	echo $(PKGNAME): Version $(VERSION)

depinfo:
	@echo -n "The package $(PKGNAME) marks [$(DEPS)] as dependencies, rpkgm does not support dependency resolution and installation yet, you will need to install them manually.\n"

clean:
	rm -rvf $(WORKDIR)