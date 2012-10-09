
all: 
	GOPATH=`pwd` go install invlist
	GOPATH=`pwd` go install treap
	GOPATH=`pwd` go install buildProgram
#	GOPATH=`pwd` go install flame

fmt:
	GOPATH=`pwd` go fmt buildProgram
	GOPATH=`pwd` go fmt invlist

test:
	GOPATH=`pwd` go test invlist

doc:
	GOPATH=`pwd` godoc -http=":8000"

clean:
	rm -rf ./bin/
	rm -rf ./pkg/
