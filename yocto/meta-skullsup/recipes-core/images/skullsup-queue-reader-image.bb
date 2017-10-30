DESCRIPTION = "Barebones image for hosting Skulls Up! queue reader"
LICENSE = "MIT"

inherit core-image

IMAGE_INSTALL += "\
    skullsup-queue-reader \
    skullsup-queue-reader-initscripts \
    skullsup-queue-reader-config \
    wpa-supplicant \
"

inherit extrausers
EXTRA_USERS_PARAMS = "\
    useradd -G dialout -r -d ${SKULLSUP_DIR} -s /bin/false skullsup; \
"
