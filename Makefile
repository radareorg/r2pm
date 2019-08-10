all: r2pm libr2pm.so

r2pm:
	go build

libr2pm.so: lib/main.go
	go build -o $@ -buildmode=c-shared ./lib

r2pm_c: c/r2pm.c libr2pm.so
	gcc -o $@ -I. -L. -lr2pm $<

clean:
	rm libr2pm.so r2pm r2pm_c r2pm.h
