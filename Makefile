# Go parameters
TEST_PKGS = $(shell go list ./...)
run:
	go run main.go

test:
	go get github.com/alecthomas/participle
	mkdir -p report
	$(foreach i,$(TEST_PKGS),go test $(i) -test.short -v -covermode=count -coverprofile=report/cover-`echo $(i) | sed 's/\//./g'`.coverprofile || exit 1;)
	rm report/*.coverprofile
