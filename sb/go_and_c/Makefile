
demo : demo.go libccode.a
	go build demo.go

ccode.o : ccode.c
	gcc -c ccode.c

libccode.a : ccode.o
	ar -crs libccode.a ccode.o
