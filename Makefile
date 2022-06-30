.PHONY: unit-test
unit-test:
	go test ./client/... -count=1

.PHONY: integration-test
integration-test:
	go test ./test/integration/... -count=1