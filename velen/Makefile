PREFIX ?= /usr
EPREFIX ?= ${PREFIX}/bin
CC ?= gcc
CFLAGS ?= -Wall -Wextra -std=c11

all: velen
.PHONY: all

%.o: %.c
	${CROSS_COMPILE}${CC} ${CFLAGS} -o $@ -c $<

velen: velen.o velen-prepare.o velen-run.o velen-destroy.o pivot_root.o detach_mount.o
	${CROSS_COMPILE}${CC} ${CFLAGS} -o $@ $^

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

detach_mount.o: detach_mount.c detach_mount.h
pivot_root.o: pivot_root.c pivot_root.h
velen.o: velen.c velen-prepare.h velen-run.h velen-destroy.h
velen-destroy.o: velen-destroy.c velen-destroy.h detach_mount.h config.h
velen-prepare.o: velen-prepare.c velen-prepare.h config.h
velen-run.o: velen-run.c velen-run.h pivot_root.h detach_mount.h config.h
