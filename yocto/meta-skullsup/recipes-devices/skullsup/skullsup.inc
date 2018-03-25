LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

S = "${WORKDIR}"

SRC_URI += "\
    file://src/skullsup/uart.go
    file://src/skullsup/hexdump.go
    file://src/skullsup/version.go
    file://src/skullsup/skullsup.go
    file://src/skullsup/frame.go
    file://src/skullsup/color.go
    file://src/skullsup/platform.go
    file://src/skullsup/psalms.go
    file://src/cmd/server/auth.go
    file://src/cmd/server/messagequeue.go
    file://src/cmd/server/server.go
    file://src/cmd/skullsup-queue-writer.go
    file://src/cmd/skullsup-queue-incantor.go
    file://src/cmd/skullsup-queue-randcolor.go
    file://src/cmd/common/common.go
    file://src/cmd/common/common_windows.go
    file://src/cmd/common/defaults/defaults_linux.go
    file://src/cmd/common/defaults/defaults_windows.go
    file://src/cmd/common/defaults/defaults.go
    file://src/cmd/common/logger/logger.go
    file://src/cmd/common/common_linux.go
    file://src/cmd/skullsup.go
    file://src/cmd/client/config.go
    file://src/cmd/client/client.go
    file://src/cmd/skullsup-queue-reader.go
    file://src/cmd/skullsup-queue-server.go
"

inherit go

GOARCH = "${TARGET_GOARCH}"
GOARM =  "${TARGET_GOARM}"

do_configure() {
    ${GO} get golang.org/x/crypto/bcrypt
    ${GO} get github.com/tarm/serial
    ${GO} get gopkg.in/alecthomas/kingpin.v2
}

do_compile() {
    ${GO} build ${S}/src/cmd/${PN}.go
}

do_install() {
    install -D -m 0755 ${B}/${PN} ${D}${bindir}/${PN}
}
