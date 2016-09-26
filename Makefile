TOOLS= golang.org/x/tools/cover
GOCOVERFILE?=	.cover.out
GOCOVERHTML?=	coverage.html

test:: $(GOCOVERFILE)

cover:: coverage_report

${GOCOVERFILE}::
	find * -type d -print0 | xargs -0 -I % sh -c "cd %; go test -v -race -cover -coverprofile=$(GOCOVERFILE)"

$(GOCOVERHTML): $(GOCOVERFILE)
	go tool cover -html=$(GOCOVERFILE) -o $(GOCOVERHTML)

coverage_report:: $(GOCOVERFILE)
	go tool cover -html=$(GOCOVERFILE)

audit_tools::
	@go get -u github.com/golang/lint/golint && echo "Installed golint:"
	@go get -u github.com/fzipp/gocyclo && echo "Installed gocyclo:"
	@go get -u github.com/remyoudompheng/go-misc/deadcode && echo "Installed deadcode:"
	@go get -u github.com/client9/misspell/cmd/misspell && echo "Installed misspell:"
	@go get -u github.com/gordonklaus/ineffassign && echo "Installed ineffassign:"

audit::
	deadcode
	go tool vet -all *.go
	go tool vet -shadow=true *.go
	golint *.go
	ineffassign .
	gocyclo -over 65 *.go
	misspell *.go

clean::
	rm -f $(GOCOVERFILE) $(GOCOVERHTML)
