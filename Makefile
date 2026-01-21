# Simple Makefile to build/run examples with the right build tags.
#
# Usage:
#   make test
#   make run-httprouter
#   make run-httprouter-security
#   make run-httprouter-typed
#   make run-httprouter-typed-security
#
#   make run-gin
#   make run-gin-security
#   make run-gin-typed
#   make run-gin-typed-security
#
#   make run-echo
#   make run-echo-security
#   make run-echo-typed
#   make run-echo-typed-security
#
#   make run-fiber
#   make run-fiber-security
#   make run-fiber-typed
#   make run-fiber-typed-security

GO ?= go

.PHONY: help test tidy fmt lint \
	run-httprouter run-httprouter-security \
	run-httprouter-typed run-httprouter-typed-security \
	run-gin run-gin-security run-gin-typed run-gin-typed-security \
	run-echo run-echo-security run-echo-typed run-echo-typed-security \
	run-fiber run-fiber-security run-fiber-typed run-fiber-typed-security

help:
	@echo "Targets:"
	@echo "  test                              - run unit tests"
	@echo "  tidy                              - go mod tidy"
	@echo "  fmt                               - gofmt all go files"
	@echo "  run-httprouter                     - run net/http example"
	@echo "  run-httprouter-security            - run net/http security example (-tags security)"
	@echo "  run-httprouter-typed               - run net/http typed example"
	@echo "  run-httprouter-typed-security      - run net/http typed security example (-tags security)"
	@echo "  run-gin                            - run gin basic example (-tags gin)"
	@echo "  run-gin-security                   - run gin basic security example (-tags gin,security)"
	@echo "  run-gin-typed                      - run gin typed example (-tags gin,typed)"
	@echo "  run-gin-typed-security             - run gin typed security example (-tags gin,typed,security)"
	@echo "  run-echo                           - run echo basic example (-tags echo)"
	@echo "  run-echo-security                  - run echo basic security example (-tags echo,security)"
	@echo "  run-echo-typed                     - run echo typed example (-tags echo,typed)"
	@echo "  run-echo-typed-security            - run echo typed security example (-tags echo,typed,security)"
	@echo "  run-fiber                          - run fiber basic example (-tags fiber)"
	@echo "  run-fiber-security                 - run fiber basic security example (-tags fiber,security)"
	@echo "  run-fiber-typed                    - run fiber typed example (-tags fiber,typed)"
	@echo "  run-fiber-typed-security           - run fiber typed security example (-tags fiber,typed,security)"


test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...
	# gofmt on files that may not be in a package list in some setups
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

tidy:
	$(GO) mod tidy

# --- net/http examples ---
run-httprouter:
	$(GO) run ./example/httprouter

run-httprouter-security:
	$(GO) run -tags security ./example/httprouter

run-httprouter-typed:
	$(GO) run ./example/httprouter_typed

run-httprouter-typed-security:
	$(GO) run -tags security ./example/httprouter_typed

# --- Gin examples ---
run-gin:
	$(GO) run -tags gin ./example/gin

run-gin-security:
	$(GO) run -tags "gin,security" ./example/gin

run-gin-typed:
	$(GO) run -tags "gin,typed" ./example/gin

run-gin-typed-security:
	$(GO) run -tags "gin,typed,security" ./example/gin

# --- Echo examples ---
run-echo:
	$(GO) run -tags echo ./example/echo

run-echo-security:
	$(GO) run -tags "echo,security" ./example/echo

run-echo-typed:
	$(GO) run -tags "echo,typed" ./example/echo

run-echo-typed-security:
	$(GO) run -tags "echo,typed,security" ./example/echo

# --- Fiber examples ---
run-fiber:
	$(GO) run -tags fiber ./example/fiber

run-fiber-security:
	$(GO) run -tags "fiber,security" ./example/fiber

run-fiber-typed:
	$(GO) run -tags "fiber,typed" ./example/fiber

run-fiber-typed-security:
	$(GO) run -tags "fiber,typed,security" ./example/fiber
