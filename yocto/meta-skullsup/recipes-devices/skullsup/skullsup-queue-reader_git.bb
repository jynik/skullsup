include skullsup.inc

# CA Certificates required for HTTPS client
RDEPENDS_${PN} += "ca-certificates"

inherit useradd
USERADD_PACKAGES = "${PN}"
USERADD_PARAM_${PN}  = "-U -G dialout -r -d ${SKULLSUP_DIR} -s /bin/false skullsup"

INFO_SUFFIX = "is not defined. Using the default value. Specify this in your local.conf if this is not desired."
FATAL_SUFFIX = "is not defined. Specify this in your local.conf if this is not desired."

_install_file() {
    local perms="$1"
    local src="$2"
    local dest="$3"

    install -m ${perms} ${src} ${dest}
    chown root:skullsup "${dest}"
}

do_install_append() {
    if [ -z "${SKULLSUP_QUEUE_READER_CONFIG}" ]; then
        bbfatal "SKULLSUP_QUEUE_READER_CONFIG ${FATAL_SUFFIX}"
    fi

    if [ -z "${SKULLSUP_QUEUE_READER_CERT}" ]; then
        bbfatal "SKULLSUP_QUEUE_READER_CERT ${FATAL_SUFFIX}"
    fi

    if [ -z "${SKULLSUP_QUEUE_READER_KEY}" ]; then
        bbfatal "SKULLSUP_QUEUE_READER_KEY ${FATAL_SUFFIX}"
    fi

    if [ -z "${SKULLSUP_DIR}" ]; then
        SKULLSUP_DIR=/opt/skullsup
        bbinfo "SKULLSUP_DIR ${INFO_SUFFIX}"
    fi

    install -d ${D}${SKULLSUP_DIR}

    mkdir ${D}${SKULLSUP_DIR}/log
    chown root:skullsup ${D}${SKULLSUP_DIR}/log
    chmod 0770 ${D}/${SKULLSUP_DIR}/log

    _install_file 0644 ${SKULLSUP_QUEUE_READER_CONFIG}  ${D}${SKULLSUP_DIR}/skullsup-queue-reader.cfg
    _install_file 0644 ${SKULLSUP_QUEUE_READER_CERT}    ${D}${SKULLSUP_DIR}/$(basename ${SKULLSUP_QUEUE_READER_CERT})
    _install_file 0644 ${SKULLSUP_QUEUE_READER_KEY}     ${D}${SKULLSUP_DIR}/$(basename ${SKULLSUP_QUEUE_READER_KEY})
}

FILES_${PN} += "${SKULLSUP_DIR}"
