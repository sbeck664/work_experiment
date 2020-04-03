
demo : demo.go libqpid.a qpid.h
	go build demo.go

qpid.o : qpid.cpp qpid.h
	g++ -c qpid.cpp

libqpid.a : qpid.o
	ar -crs libqpid.a qpid.o
