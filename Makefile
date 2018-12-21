# Go parameters
TEST_PKGS = $(shell go list ./...)
run:
	go run main.go

test:
	go get github.com/alecthomas/participle
	go get github.com/satori/go.uuid
	go get github.com/golang/protobuf/proto
	go get github.com/soheilhy/cmux
	go get google.golang.org/grpc
	go get github.com/brg-liuwei/godnf
	go get github.com/kr/pretty
	mkdir -p report
	$(foreach i,$(TEST_PKGS),go test $(i) -test.short -v -covermode=count -coverprofile=report/cover-`echo $(i) | sed 's/\//./g'`.coverprofile || exit 1;)
	rm report/*.coverprofile
