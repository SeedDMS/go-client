PKGNAME=seeddms-client
VERSION=0.0.2

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOGET=$(GOCMD) get

default: build

# Setting the version isn't working, but I leave it here for later
# debugging
build: test
	mkdir -p bin
	$(GOBUILD) -v -o bin/seeddms-client -ldflags "-X cmd.version=${VERSION}" main.go 

clean:
	$(GOCLEAN)
	rm -rf bin

test:
	$(GOVET) ./...
	$(GOTEST) -v ./...

run:
	bin/seeddms-client --profile steinmann_dms ls

dist: clean
	rm -rf ${PKGNAME}-${VERSION}
	mkdir ${PKGNAME}-${VERSION}
	cp -r cmd config assets *.go resources Makefile README.md seeddms-client.yaml go.mod go.sum ${PKGNAME}-${VERSION}
	tar czvf ${PKGNAME}-${VERSION}.tar.gz ${PKGNAME}-${VERSION}
	rm -rf ${PKGNAME}-${VERSION}

debian: dist
	mv ${PKGNAME}-${VERSION}.tar.gz ../${PKGNAME}_${VERSION}.orig.tar.gz
	debuild

.PHONY: build test debian
