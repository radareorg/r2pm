ifeq ($(OS),Windows_NT)
    LIB_EXT := .dll
else
    LIB_EXT := .so
endif

LIB := libr2pm${LIB_EXT}

all: r2pm r2pm_c ${LIB}

.PHONY: tests integration-tests

integration-tests:
	go test -v -tags=integration ./...

tests:
	go test ./...

r2pm: $(wildcard internal/**/*.go pkg/**/*.go main.go)
	go build

${LIB}: $(wildcard internal/**/*.go lib/*.go pkg/**/*.go)
	go build -o $@ -buildmode=c-shared ./lib

r2pm_c: c/r2pm.c ${LIB}
	${CC} -Wall -o $@ -I. -L. $< -lr2pm

clean:
	rm -f ${LIB} libr2pm.h r2pm r2pm_c
