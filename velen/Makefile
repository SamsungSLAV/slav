PREFIX ?= /usr
EPREFIX ?= ${PREFIX}/bin
CC ?= gcc
CFLAGS ?= -MMD -MP -Wall -Wextra -std=c11

SOURCES = velen.c velen-prepare.c velen-run.c velen-destroy.c pivot_root.c detach_mount.c
OBJECTS = $(SOURCES:.c=.o)

all: velen
.PHONY: all

%.o: %.c
	${CROSS_COMPILE}${CC} ${CFLAGS} -o $@ -c $<

velen: ${OBJECTS}
	${CROSS_COMPILE}${CC} ${CFLAGS} -o $@ $^

install: velen
	mkdir -p ${DESTDIR}${EPREFIX}
	install -m755 velen ${DESTDIR}${EPREFIX}/velen
	ln -fsv velen ${DESTDIR}${EPREFIX}/velen-prepare
	ln -fsv velen ${DESTDIR}${EPREFIX}/velen-run
	ln -fsv velen ${DESTDIR}${EPREFIX}/velen-destroy
.PHONY: install

clean:
	-rm *.o *.d velen
.PHONY: clean

-include $(SOURCES:.c=.d)
