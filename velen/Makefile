PREFIX ?= /usr
EPREFIX ?= ${PREFIX}/bin
CC ?= gcc

all: velen
.PHONY: all

velen: velen.c velen-prepare.c config.h velen-run.c pivot_root.c velen-destroy.c
	${CROSS_COMPILE}${CC} -o $@ velen.c

install: velen
	mkdir -p ${DESTDIR}${EPREFIX}
	install -m755 velen ${DESTDIR}${EPREFIX}/velen
	ln -fsv velen ${DESTDIR}${EPREFIX}/velen-prepare
	ln -fsv velen ${DESTDIR}${EPREFIX}/velen-run
	ln -fsv velen ${DESTDIR}${EPREFIX}/velen-destroy
.PHONY: install
