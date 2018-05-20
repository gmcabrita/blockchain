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
	gometalinter --vendor --linter="prealloc:prealloc:^(?P<path>.*?\\.go):(?P<line>\\d+)\\s*(?P<message>.*)$$" --enable=prealloc --disable=vetshadow --line-length=120 ./...;

# installs test dependencies
testdeps:
	go get -u -v github.com/mattn/goveralls;

# tests the project
test:
	go test -race -v -covermode=atomic -coverprofile=coverage.out ./...;
