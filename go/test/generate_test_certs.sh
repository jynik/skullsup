#!/bin/sh
#
# Quick and dirty script to generate some test keys and certificates
################################################################################

THIS_DIR=$(realpath $(dirname $0))

KEYS_DIR="${THIS_DIR}/keys"
CERTS_DIR="${THIS_DIR}/certs"
CFG_DIR="${THIS_DIR}/configs"

gen_test_client_cert()
{
	clientno=$1
	client="client${clientno}"
	echo "Generating client certificate for $client"
	openssl genrsa -out "${KEYS_DIR}/${client}.key" 2048 >/dev/null 2>&1

	openssl req -new -key "${KEYS_DIR}/${client}.key" -out "${CERTS_DIR}/${client}".csr >/dev/null 2>&1 <<EOF
US
Hell
MI
SkullsUp!
Chaos Division
$client



EOF
	openssl x509 -req -in "${CERTS_DIR}/${client}.csr" -CA "${CERTS_DIR}/ca.pem" -CAkey "${KEYS_DIR}/ca.key" -set_serial ${clientno} -out "${CERTS_DIR}/${client}.pem" >/dev/null 2>&1
}

echo 
echo This script is intended for creating test keys and certificates.
echo The resulting files are not suitable for deployments.
echo  Press enter to acknowledge this and continue...
echo ----------------------------------------------------------------
read cya_disclaimer

mkdir -p ${KEYS_DIR} ${CERTS_DIR}

echo Creating test server certificate...
yes "" | openssl req -nodes -new -x509 -keyout "${KEYS_DIR}/skullsup-queue-server.key" -out "${CERTS_DIR}/skullsup-queue-server.pem" >/dev/null 2>&1

echo Creating test CA key
openssl genrsa -out "${KEYS_DIR}/ca.key" 2048 >/dev/null 2>&1

echo Creating test CA certificate
yes "" | openssl req -new -x509 -key "${KEYS_DIR}/ca.key" -out "${CERTS_DIR}/ca.pem" >/dev/null 2>&1

echo Creating test client certificates
for n in 1 2 3 4; do 
	gen_test_client_cert $n
done
