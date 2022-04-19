.POSIX:

bindir = ./bin

all: mazegenerator mazesolver

clean:
	rm -f $(bindir)/mazegenerator $(bindir)/mazesolver
	rmdir $(bindir)

mazegenerator: mazegenerator/mazegenerator.go
	go build -o $(bindir)/mazegenerator ./mazegenerator

mazesolver: mazesolver.go
	go build -o $(bindir)/mazesolver .

.PHONY: all clean mazegenerator mazesolver
