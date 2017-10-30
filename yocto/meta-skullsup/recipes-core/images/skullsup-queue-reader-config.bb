DESCRIPTION = "Installs SkullsUp! Que Reader configuration file"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

INFO_SUFFIX = "is not defined. Using the default value. Specify this in your local.conf if this is not desired." 

do_install() {
    if [ -z "${SKULLSUP_QUEUE_READER_CONFIG}" ]; then
        bbfatal "SKULLSUP_QUEUE_READER_CONFIG is not defined. Specify this as a full path in your local.conf."
    fi

    if [ -z "${SKULLSUP_DIR}" ]; then
        SKULLSUP_DIR=/opt/skullsup
        bbfatal "SKULLSUP_DIR is not defined in local.conf."
    fi

    install -d ${D}${SKULLSUP_DIR} 
    install -m 644 ${SKULLSUP_QUEUE_READER_CONFIG} ${D}${SKULLSUP_DIR}
}

FILES_${PN} = "${SKULLSUP_DIR}"
