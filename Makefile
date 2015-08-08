.PHONY: all install test vet

all: test vet

test:
	go test

vet:
	go vet
	#golint .
