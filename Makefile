# Corvint Tasks verification gate.
# Fails closed until Go source exists so an empty tree can never read as gate-green.

GO ?= go
GOENV := GOTOOLCHAIN=local
WANT_GO := go1.27.1

.PHONY: verify toolchain sources fmt test vet

verify: toolchain sources fmt test vet

toolchain:
	@have="$$($(GOENV) $(GO) env GOVERSION)"; \
	if [ "$$have" != "$(WANT_GO)" ]; then \
		echo "toolchain gate: want $(WANT_GO), have $$have"; exit 1; \
	fi

sources:
	@if [ -z "$$(find . -name '*.go' -not -path './.git/*' -print -quit)" ]; then \
		echo "verify gate: NOT_RUN (no Go source exists yet; see docs/ROADMAP.md)"; exit 1; \
	fi

# gofmt exits non-zero on a syntax error (reported on stderr, nothing on
# stdout) and zero with a file list when formatting differs; both fail.
fmt: sources
	@out="$$(gofmt -l . 2>&1)"; status=$$?; \
	if [ "$$status" -ne 0 ]; then echo "gofmt failed (exit $$status):"; echo "$$out"; exit "$$status"; fi; \
	if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi

test: sources
	$(GOENV) $(GO) test -count=1 ./...

vet: sources
	$(GOENV) $(GO) vet ./...
