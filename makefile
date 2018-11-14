GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
PLUGINDIR= plugins
BINARY_NAME= start

PLUGINSOURCES:= $(shell find $(PLUGINDIR) -name '*.go')

#all: test build

build:
	$(foreach element, $(PLUGINSOURCES), $(GOBUILD) -buildmode=plugin -o $(patsubst %.go,%.so, $(element)) $(element);)
	$(GOBUILD) -o $(BINARY_NAME) -v

#test: $(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	$(foreach plugin, $(shell find $(PLUGINDIR) -name '*.so') , rm -f $(plugin))

#run: build
#	./$(BINARY_NAME)
