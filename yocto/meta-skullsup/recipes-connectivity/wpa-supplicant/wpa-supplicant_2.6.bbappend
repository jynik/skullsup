FILESEXTRAPATHS_prepend := "${THISDIR}/files:"

do_configure_prepend() {
    if [ ! -z "${SKULLSUP_WIFI_SSID}" ]; then
        sed -i -e "s/#ssid=.*$/ssid=\"${SKULLSUP_WIFI_SSID}\"/" ${WORKDIR}/wpa_supplicant.conf-sane
    fi

    # We expect the psk as the hash here, 
    if [ ! -z "${SKULLSUP_WPA_PSK}" ]; then
        sed -i -e "s/#psk=.*$/psk=${SKULLSUP_WPA_PSK}/" ${WORKDIR}/wpa_supplicant.conf-sane
    fi
}

