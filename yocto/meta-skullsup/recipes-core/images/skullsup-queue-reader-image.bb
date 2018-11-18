DESCRIPTION = "Barebones image for hosting Skulls Up! queue reader"
LICENSE = "MIT"

inherit core-image

IMAGE_INSTALL += "\
    init-iptables \
    ntp \
    rng-tools \
    skullsup-queue-reader \
    skullsup-queue-reader-initscripts \
    wpa-supplicant \
"

