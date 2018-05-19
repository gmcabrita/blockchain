# installs project dependencies
deps:
	go get -u -v github.com/golang/dep/cmd/dep;
	dep ensure;

# installs lint dependencies
lintdeps:
	go get -u -v github.com/alexkohler/prealloc \
		github.com/alecthomas/gometalinter;
	gometalinter --install;

# lints the project
lint:
	gometalinter --linter="prealloc:prealloc:^(?P<path>.*?\\.go):(?P<line>\\d+)\\s*(?P<message>.*)$$" --enable-all --enable=prealloc ./...;

# installs test dependencies
testdeps:
	go get -u -v github.com/mattn/goveralls;

# tests the project
test:
	go test -race -v -covermode=atomic -coverprofile=coverage.out ./...;