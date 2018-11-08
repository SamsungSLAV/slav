PREFIX ?= /usr
EPREFIX ?= ${PREFIX}/bin
CC ?= gcc
CFLAGS ?= -Wall -Wextra -std=c11

all: velen
.PHONY: all

%.o: %.c
	${CROSS_COMPILE}${CC} -o $@ -c $<

velen: velen.o velen-prepare.o velen-run.o pivot_root.o velen-destroy.o
	${CROSS_COMPILE}${CC} -o $@ $^

install: velen
	mkdir -p ${DESTDIR}${EPREFIX}
	install -m755 velen ${DESTDIR}${EPREFIX}/velen
	ln -fsv velen ${DESTDIR}${EPREFIX}/velen-prepare
	ln -fsv velen ${DESTDIR}${EPREFIX}/velen-run
	ln -fsv velen ${DESTDIR}${EPREFIX}/velen-destroy
.PHONY: install

clean:
	-rm *.o velen
.PHONY: clean

pivot_root.o: pivot_root.c pivot_root.h
velen.o: velen.c velen-prepare.h velen-run.h velen-destroy.h
velen-destroy.o: velen-destroy.c velen-destroy.h config.h
velen-prepare.o: velen-prepare.c velen-prepare.h config.h
velen-run.o: velen-run.c velen-run.h pivot_root.h config.h
