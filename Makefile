GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test

install:
	$(GOINSTALL)

clean:
	$(GOCLEAN) -n -i -x
	rm -f $(GOPATH)/bin/logcolor
	rm -f logcolor

test:
	$(GOTEST) -v -cover
