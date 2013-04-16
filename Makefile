MGO_BIN=$(GOPATH)/bin/mgo

all:	$(MGO_BIN)

$(MGO_BIN): */*.go
	cd mgo ; go install -v

# XXX This should ideally go test in all dirs that contain _test.go files,
# but make isn't cooperating
test:
	(cd mgo_connector_redis ; go test -v)

fmt:
	gofmt -s -l -w .
