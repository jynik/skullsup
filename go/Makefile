SKULLSUP_LIB_SRC := $(wildcard src/skullsup/*.go)
SKULLSUP_CMD_SRC := $(wildcard src/cmd/common/*.go) \
					$(wildcard src/cmd/common/defaults/*.go) \
					$(wildcard src/cmd/client/*.go) \
					$(wildcard src/cmd/server/*.go) \
					$(wildcard src/cmd/*.go)

SKULLSUP := src/cmd/skullsup.go
SKULLSUP_SERVER := src/cmd/skullsup-queue-server.go
SKULLSUP_READER := src/cmd/skullsup-queue-reader.go
SKULLSUP_WRITER := src/cmd/skullsup-queue-writer.go

SKULLSUP_ALL := $(SKULLSUP_LIB_SRC) $(SKULLSUP_CMD_SRC) $(SKULLSUP) \
				$(SKULLSUP_SERVER) $(SKULLSUP_WRITER)


all: skullsup skullsup-queue-server skullsup-queue-reader skullsup-queue-writer

deps:
	go get golang.org/x/crypto/bcrypt
	go get github.com/tarm/serial
	go get gopkg.in/alecthomas/kingpin.v2

format:
	gofmt -w $(SKULLSUP_ALL)

skullsup: $(SKULLSUP_LIB_SRC) $(SKULLSUP_CMD_SRC) $(SKULLSUP)
	go build $(SKULLSUP)

skullsup-queue-server: $(SKULLSUP_LIB_SRC) $(SKULLSUP_CMD_SRC) $(SKULLSUP_SERVER)
	go build $(SKULLSUP_SERVER)

skullsup-queue-reader: $(SKULLSUP_LIB_SRC) $(SKULLSUP_CMD_SRC) $(SKULLSUP_READER)
	go build $(SKULLSUP_READER)

skullsup-queue-writer: $(SKULLSUP_LIB_SRC) $(SKULLSUP_CMD_SRC) $(SKULLSUP_WRITER)
	go build $(SKULLSUP_WRITER)

clean:
	rm -f skullsup skullsup-queue-server skullsup-queue-writer skullsup-queue-reader

.PHONY: doc format clean
