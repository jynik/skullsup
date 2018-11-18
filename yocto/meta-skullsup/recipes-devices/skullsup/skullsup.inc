LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

S = "${WORKDIR}"

inherit go

GOARCH = "${TARGET_GOARCH}"
GOARM =  "${TARGET_GOARM}"

do_configure() {
    ${GO} get github.com/tarm/serial
    ${GO} get github.com/jynik/skullsup/go/src/cmd/${PN}
}

do_compile() {
    ${GO} build ./src/github.com/jynik/skullsup/go/src/cmd/${PN}
}

do_install() {
    install -D -m 0755 ${B}/${PN} ${D}${bindir}/${PN}
}
