MGO_BIN=$(GOPATH)/bin/mgo

all:	$(MGO_BIN)

$(MGO_BIN): mgo/mgo.go *.go
	cd mgo ; go install -v

test:
	go test -v

fmt:
	gofmt -s -l -w .
