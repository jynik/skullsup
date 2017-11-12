DESCRIPTION = "Startup script for loading iptables rules"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

RDEPENDS_${PN} = "iptables"

S = "${WORKDIR}"

SRC_URI = "\
    file://iptables \
    file://iptables.conf \
"

do_install() {
    install -d ${D}${sysconfdir}/init.d
    install -d ${D}${sysconfdir}/rcS.d
    install -d ${D}${sysconfdir}/rc1.d
    install -d ${D}${sysconfdir}/rc2.d
    install -d ${D}${sysconfdir}/rc3.d
    install -d ${D}${sysconfdir}/rc4.d
    install -d ${D}${sysconfdir}/rc5.d

    install -m 0755 ${S}/iptables ${D}${sysconfdir}/init.d/
    install -m 0700 ${S}/iptables.conf ${D}/${sysconfdir}/

    ln -sf ../init.d/iptables ${D}${sysconfdir}/rcS.d/S01iptables
    ln -sf ../init.d/iptables ${D}${sysconfdir}/rc2.d/S01iptables
    ln -sf ../init.d/iptables ${D}${sysconfdir}/rc3.d/S01iptables
    ln -sf ../init.d/iptables ${D}${sysconfdir}/rc4.d/S01iptables
    ln -sf ../init.d/iptables ${D}${sysconfdir}/rc5.d/S01iptables
}
