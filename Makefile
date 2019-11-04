DEBUGFLAGS = -g -ggdb -O2
ifeq ($(DEBUG), 1)
	DEBUGFLAGS = -g -ggdb -O0
endif

CFLAGS = -Wall -Wno-unused-function ${DEBUGFLAGS} -fPIC -std=c11 -D_GNU_SOURCE
CC := $(shell sh -c 'type $(CC) >/dev/null 2>/dev/null && echo $(CC) || echo gcc')

# Compile flags for linux / osx
uname_S := $(shell sh -c 'uname -s 2>/dev/null || echo not')
ifeq ($(uname_S), Linux)
	SHOBJ_CFLAGS ?= -fno-common -g -ggdb
	SHOBJ_LDFLAGS ?= -shared -Bsymbolic -Bsymbolic-functions
else
	CFLAGS += -mmacosx-version-min=10.6
	SHOBJ_CFLAGS ?= -dynamic -fno-common -g -ggdb
	SHOBJ_LDFLAGS ?= -dylib -exported_symbol _RedisModule_OnLoad -macosx_version_min 10.6
endif


CC_SOURCES = $(wildcard *.c)
CC_OBJECTS = $(patsubst %.c, %.o, ${CC_SOURCES})

all: libredis_ratelimit.so

libredis_ratelimit.so: ${CC_OBJECTS} libs/libratelimit.a
	${LD} -o $@ ${CC_OBJECTS} $(SHOBJ_LDFLAGS) -lc -lpthread -Llibs -lratelimit

libs/libratelimit.a:
	make -C libs

clean:
	make clean -C libs
	rm -rvf *.xo *.so *.o *.a
