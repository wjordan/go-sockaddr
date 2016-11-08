BIN:=sockaddr
SRCS:=$(shell find . -name '*.go' ! -path '*/vendor/*')

.DEFAULT_GOAL := dev

.PHONY: dev
dev: $(BIN)
	@install $(BIN) ${GOPATH}/bin/

$(BIN): $(SRCS)
	go build -o $@

.PHONY: clean
clean::
	rm -f $(BIN) regression/*.diff

.PHONY: install
install:: $(BIN)
	install sockaddr ${GOPATH}/bin/

.PHONY: test
test::	$(BIN)
	@make -C regression
