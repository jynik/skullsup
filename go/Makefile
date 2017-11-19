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
SKULLSUP_INCANTOR := src/cmd/skullsup-queue-incantor.go
SKULLSUP_RANDCOLOR := src/cmd/skullsup-queue-randcolor.go

COMMANDS := ${SKULLSUP} \
			${SKULLSUP_SERVER} \
			${SKULLSUP_READER} \
			${SKULLSUP_WRITER} \
			${SKULLSUP_INCANTOR} \
			${SKULLSUP_RANDCOLOR}

BINARIES := $(notdir $(COMMANDS:.go=))

SKULLSUP_ALL := $(SKULLSUP_LIB_SRC) $(SKULLSUP_CMD_SRC) $(SKULLSUP) \
				$(SKULLSUP_SERVER) $(SKULLSUP_WRITER)

all: ${BINARIES}

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

skullsup-queue-incantor: ${SKULLSUP_LIB_SRC} ${SKULLSUP_CMD_SRC} ${SKULLSUP_INCANTOR}
	go build $(SKULLSUP_INCANTOR)

skullsup-queue-randcolor: ${SKULLSUP_LIB_SRC} ${SKULLSUP_CMD_SRC} ${SKULLSUP_RANDCOLOR}
	go build $(SKULLSUP_RANDCOLOR)

clean:
	rm -f ${BINARIES}

.PHONY: doc format clean
