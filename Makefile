all: r2pm r2pm_c libr2pm.so

.PHONY: test

test:
	go test ./...

r2pm: $(wildcard internal/**/*.go pkg/**/*.go main.go)
	go build

libr2pm.so: $(wildcard internal/**/*.go lib/*.go pkg/**/*.go)
	go build -o $@ -buildmode=c-shared ./lib

r2pm_c: c/r2pm.c libr2pm.so
	gcc -Wall -o $@ -I. -L. $< -lr2pm

clean:
	rm libr2pm.so libr2pm.h r2pm r2pm_c
