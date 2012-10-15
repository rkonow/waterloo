
all: 
	GOPATH=`pwd` go install invlist
	GOPATH=`pwd` go install treap
	GOPATH=`pwd` go install hashlist
	GOPATH=`pwd` go install svs
	GOPATH=`pwd` go install buildProgram

fmt:
	GOPATH=`pwd` go fmt buildProgram
	GOPATH=`pwd` go fmt invlist
	GOPATH=`pwd` go fmt hashlist
	GOPATH=`pwd` go fmt svs
	GOPATH=`pwd` go fmt treap

test:
	GOPATH=`pwd` go test invlist

doc:
	GOPATH=`pwd` godoc -http=":8000"

clean:
	rm -rf ./bin/
	rm -rf ./pkg/
