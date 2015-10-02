# Go parameters
GOCMD ?=     go
GOBUILD ?=   $(GOCMD) build
GOCLEAN ?=   $(GOCMD) clean
GOINSTALL ?= $(GOCMD) install
GOTEST ?=    $(GOCMD) test
GOFMT ?=     gofmt -w

NAME = pathwar-cli
DEPS = .deps
SRC = .
PACKAGES = $(shell find ./pkg/* -type d)

BUILD_LIST = $(foreach int, $(SRC), $(int)_build)
CLEAN_LIST = $(foreach int, $(SRC) $(PACKAGES), $(int)_clean)
INSTALL_LIST = $(foreach int, $(SRC), $(int)_install)
IREF_LIST = $(foreach int, $(SRC) $(PACKAGES), $(int)_iref)
TEST_LIST = $(foreach int, $(SRC) $(PACKAGES), $(int)_test)
FMT_LIST = $(foreach int, $(SRC) $(PACKAGES), $(int)_fmt)

.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(BUILD_LIST) $(IREF_LIST)

all: build
build: $(DEPS) $(BUILD_LIST)
clean: $(CLEAN_LIST)
install: $(INSTALL_LIST)
test: $(TEST_LIST)
iref: $(IREF_LIST)
fmt: $(FMT_LIST)

$(DEPS) :
	godep 2> /dev/null || go get github.com/tools/godep
	godep restore
	touch $@

$(BUILD_LIST): %_build: %_fmt %_iref
	$(GOBUILD) -o $(NAME) ./$*
	go tool vet -all=true $(SRC)
	go tool vet -all=true $(PACKAGES)
$(CLEAN_LIST): %_clean:
	$(GOCLEAN) ./$*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) ./$*
$(IREF_LIST): %_iref:
	$(GOTEST) -i ./$*
$(TEST_LIST): %_test:
	$(GOTEST) ./$*
$(FMT_LIST): %_fmt:
	$(GOFMT) ./$*
