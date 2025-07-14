tidy:
	go mod tidy

lint:
	golangci-lint cache clean
	golangci-lint run -v

.PHONY: tidy lint