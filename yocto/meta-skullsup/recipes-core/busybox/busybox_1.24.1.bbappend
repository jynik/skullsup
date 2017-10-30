# Enable setuidgid command for use in initscripts
do_configure_prepend() {
    sed -iorig 's/# CONFIG_SETUIDGID is not set/CONFIG_SETUIDGID=y/' ${WORKDIR}/defconfig
}
