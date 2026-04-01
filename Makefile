build:
	CGO_ENABLED=0 go build -o prairie ./cmd/prairie/

run: build
	./prairie

test:
	go test ./...

clean:
	rm -f prairie

.PHONY: build run test clean
