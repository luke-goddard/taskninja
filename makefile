
# Watch and run tests
test:
	gotestsum --format testname --watch ./... -v

# Build and run the development build
run:
	go run cmd/taskninja.go


# Will hot reload the development environment if air is installed
# Air is a Go live reload tool that is similar to nodemon in Node.js
# NOTE: TUI acts weird as air takes over the process
# https://github.com/air-verse/air
run-watch:
	air

# View the sqlite database using sqlitebrowser
browse:
	sqlitebrowser $(shell cat $$HOME/.config/taskninja/config.yaml | grep path | awk '{print $$2}' | xargs)&

log:
	cat /tmp/taskninja.log

log-watch:
	tail -f /tmp/taskninja.log
