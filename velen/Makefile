PREFIX ?= /usr
EPREFIX ?= ${PREFIX}/bin
CC ?= gcc
CFLAGS ?= -MMD -MP -Wall -Wextra -std=c11

SOURCES = velen.c velen-prepare.c velen-run.c velen-destroy.c pivot_root.c detach_mount.c velen-shell.c
OBJECTS = $(SOURCES:.c=.o)

all: velen velen-shell
.PHONY: all

%.o: %.c
	${CROSS_COMPILE}${CC} ${CFLAGS} -o $@ -c $<

velen: ${OBJECTS}
	${CROSS_COMPILE}${CC} ${CFLAGS} -o $@ $^

velen-shell: velen-shell.c velen-run.o detach_mount.o pivot_root.o
	${CROSS_COMPILE}${CC} ${CFLAGS} -DVELEN_SHELL_BUILD -o $@ $^

install: velen
	mkdir -p ${DESTDIR}${EPREFIX}
	install -m755 velen ${DESTDIR}${EPREFIX}/velen
	install -m755 velen-shell ${DESTDIR}${EPREFIX}/velen-shell
	chmod u+s ${DESTDIR}${EPREFIX}/velen-shell
	ln -fsv velen ${DESTDIR}${EPREFIX}/velen-prepare
	ln -fsv velen ${DESTDIR}${EPREFIX}/velen-run
	ln -fsv velen ${DESTDIR}${EPREFIX}/velen-destroy
.PHONY: install

clean:
	-rm *.o *.d velen velen-shell
.PHONY: clean

-include $(SOURCES:.c=.d)
