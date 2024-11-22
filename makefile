
# Watch and run tests
test:
	gotestsum --format testname --watch ./... -v

# Build and run the development build
run:
	go run cmd/taskninja.go
