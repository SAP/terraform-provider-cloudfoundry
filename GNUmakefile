default: test

# Run acceptance tests
.PHONY: test
test:
	go test -cover  -count=1 ./... -v $(TESTARGS) -timeout 10m
