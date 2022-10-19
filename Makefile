DESTDIR ?=
IBUS_INSTALL_DIR ?= /usr/share
SHIN_DIR = $(DESTDIR)$(IBUS_INSTALL_DIR)/ibus-shin
IBUS_COMPONENT_DIR = $(DESTDIR)$(IBUS_INSTALL_DIR)/ibus/component

all: build shin.xml

build:
	go build

shin.xml: shin.xml.in
	sed 's:$$(SHIN_DIR):$(SHIN_DIR):g' shin.xml.in > shin.xml

install:
	mkdir -p '$(SHIN_DIR)'
	mkdir -p '$(IBUS_COMPONENT_DIR)'
	cp shin '$(SHIN_DIR)/'
	cp shin.xml '$(IBUS_COMPONENT_DIR)/'
