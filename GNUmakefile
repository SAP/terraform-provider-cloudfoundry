default: test

# Run acceptance tests
.PHONY: test
test:
	go test -count=1 ./... -v $(TESTARGS) -timeout 10m
