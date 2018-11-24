# Go parameters
TEST_PKGS = $(shell go list ./...)
run:
	go run main.go

test:
	mkdir -p report
	$(foreach i,$(TEST_PKGS),go test $(i) -test.short -v -covermode=count -coverprofile=report/cover-`echo $(i) | sed 's/\//./g'`.coverprofile || exit 1;)
