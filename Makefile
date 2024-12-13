.POSIX:
.SUFFIXES:

GO=go
GIT=git
RM=rm

all: pre-commit

pre-commit: tidy fmt lint vulnerabilities test # Runs all pre-commit checks.

commit: pre-commit # Commits the changes to the repository.
	$(GIT) commit -s

push: # Pushes the changes to the repository.
	$(GIT) push origin trunk

doc: # Serves the documentation locally.
	$(GO) run golang.org/x/tools/cmd/godoc@latest -http=localhost:1967

tidy: # Updates the go.mod and go.sum files.
	$(GO) mod tidy

fmt: # Formats Go source files in this repository.
	$(GO) run mvdan.cc/gofumpt@latest -e -extra -w .

lint: # Runs golangci-lint using the config at the root of the repository.
	$(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...

vulnerabilities: # Analyzes the codebase and looks for vulnerabilities affecting it.
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

test: # Runs unit tests.
	$(GO) test -cover -race -vet all -mod readonly ./...

test/coverage: # Generates a coverage profile and open it in a browser.
	$(GO) test -coverprofile cover.out -race -vet all -mod readonly ./...
	$(GO) tool cover -html=cover.out

.PHONY: all pre-commit commit push doc tidy fmt lint vulnerabilities test test/coverage
