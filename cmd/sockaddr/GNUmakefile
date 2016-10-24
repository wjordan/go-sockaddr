BIN:=sockaddr
SRCS:=$(shell find . -name '*.go' ! -path '*/vendor/*')

.DEFAULT_GOAL := $(BIN)

$(BIN): $(SRCS)
	go build -o $@

clean::
	rm -f $(BIN)

test::	$(BIN)
	@make -C regression
