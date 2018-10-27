DESCRIPTION = "Startup script for skullsup-queue-reader daemon"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

RDEPENDS_${PN} = "daemonize"

S = "${WORKDIR}"

DAEMON = "skullsup-queue-reader-daemon"
DAEMON_ENV = "${DAEMON}.env"
DAEMON_STARTUP = "${DAEMON}.startup"

SRC_URI = "\
    file://${DAEMON} \
    file://${DAEMON_ENV} \
    file://${DAEMON_STARTUP} \
"

INFO_SUFFIX = "is not defined. Using the default value. Specify this in your local.conf if this is not desired." 

do_configure() {
    if [ ! -z "${SKULLSUP_DIR}" ]; then
        sed -ie "s@^SKULLSUP_DIR=.*@SKULLSUP_DIR=${SKULLSUP_DIR}@g" ${S}/${DAEMON_ENV}
    else
        bbnote "SKULLSUP_DIR ${INFO_SUFFIX}"
    fi

    if [ -z "${SKULLSUP_DEVICE}" ]; then
        bbnote "SKULLSUP_DEVICE ${INFO_SUFFIX}"
    else
        sed -ie "s/^SKULLSUP_DEVICE=.*/SKULLSUP_DEVICE=${SKULLSUP_DEVICE}/g" ${S}/${DAEMON_ENV}
    fi

    if [ -z "${SKULLSUP_REMOTE_HOST}" ]; then
        bbfatal "SKULLSUP_REMOTE_HOST is not defined. Specify this in your local.conf"
    else
        sed -ie "s/^SKULLSUP_REMOTE_HOST=.*/SKULLSUP_REMOTE_HOST=${SKULLSUP_REMOTE_HOST}/g" ${S}/${DAEMON_ENV}
    fi

    if [ -z "${SKULLSUP_REMOTE_PORT}" ]; then
        bbnote "SKULLSUP_REMOTE_PORT ${INFO_SUFFIX}"
    else
        sed -ie "s/^SKULLSUP_REMOTE_PORT=.*/SKULLSUP_REMOTE_PORT=${SKULLSUP_REMOTE_PORT}/g" ${S}/${DAEMON_ENV}
    fi

    if [ -z "${SKULLSUP_POLL_PERIOD}" ]; then
        bbnote "SKULLSUP_POLL_PERIOD ${INFO_SUFFIX}"
    else
        sed -ie "s/^SKULLSUP_POLL_PERIOD=.*/SKULLSUP_POLL_PERIOD=${SKULLSUP_POLL_PERIOD}/g" ${S}/${DAEMON_ENV}
    fi
}

do_install() {
    install -d ${D}${sysconfdir}/init.d
    install -d ${D}${sysconfdir}/rcS.d
    install -d ${D}${sysconfdir}/rc1.d
    install -d ${D}${sysconfdir}/rc2.d
    install -d ${D}${sysconfdir}/rc3.d
    install -d ${D}${sysconfdir}/rc4.d
    install -d ${D}${sysconfdir}/rc5.d

    install -m 0755 ${S}/${DAEMON} ${D}${sysconfdir}/init.d/
    install -m 0755 ${S}/${DAEMON_ENV} ${D}${sysconfdir}/init.d/
    install -m 0755 ${S}/${DAEMON_STARTUP} ${D}${sysconfdir}/init.d/


    ln -sf ../init.d/${DAEMON_STARTUP} ${D}${sysconfdir}/rcS.d/S90-${DAEMON_STARTUP}

    ln -sf ../init.d/${DAEMON} ${D}${sysconfdir}/rc5.d/S90-${DAEMON}
    ln -sf ../init.d/${DAEMON} ${D}${sysconfdir}/rc1.d/K90-${DAEMON}
    ln -sf ../init.d/${DAEMON} ${D}${sysconfdir}/rc2.d/K90-${DAEMON}
    ln -sf ../init.d/${DAEMON} ${D}${sysconfdir}/rc3.d/K90-${DAEMON}
    ln -sf ../init.d/${DAEMON} ${D}${sysconfdir}/rc4.d/K90-${DAEMON}
    ln -sf ../init.d/${DAEMON} ${D}${sysconfdir}/rc5.d/K90-${DAEMON}
}
