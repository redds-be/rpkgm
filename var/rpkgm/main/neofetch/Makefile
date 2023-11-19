PKGNAME = neofetch
ARCHIVE = https://github.com/dylanaraps/neofetch/archive/refs/tags/7.1.0.tar.gz
WORKDIR = /tmp/$(PKGNAME)-build
DOWNLOADER = /usr/bin/wget -O $(WORKDIR)/$(PKGNAME)-archive.tar.gz
DESTDIR = /usr/bin
MANDIR = /usr/share/man/man1

all: dl extract install clean

dl:
	mkdir -p $(WORKDIR)
	$(DOWNLOADER) $(ARCHIVE)

extract:
	tar -xf $(WORKDIR)/$(PKGNAME)-archive.tar.gz -C $(WORKDIR) --strip-components=1

install:
	mkdir -p $(DESTDIR)
	mkdir -p $(MANDIR)
	install -o root -g root -m 0755 -p $(WORKDIR)/neofetch $(DESTDIR)/$(PKGNAME)
	install -o root -g root -m 0644 -p $(WORKDIR)/neofetch.1 $(MANDIR)

uninstall:
	rm -rf $(DESTDIR)/$(PKGNAME)
	rm -rf $(MANDIR)/neofetch.1*

clean:
	rm -rf $(WORKDIR)