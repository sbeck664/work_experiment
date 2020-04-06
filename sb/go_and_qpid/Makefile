
demo : demo.go qpid/qpid.go libqpid.a qpid.h
	go build demo.go

qpid.o : qpid.cpp qpid.h
	g++ -std=c++17 -Wall -Wextra -Wpedantic  -Werror -Weffc++ -Wold-style-cast -c qpid.cpp

libqpid.a : qpid.o
	ar -crs libqpid.a qpid.o
