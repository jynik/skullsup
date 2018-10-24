#!/bin/bash
#
# Run color and incant commands until forcefully stopped
################################################################################
set -e

if [ $# -ne 1 ]; then
	echo "Usage: $0 <interface>" >&2
	exit 1
fi

this_dir=$(dirname $(realpath "${0}"))
interface="$1"
i=0


while true; do
	echo -ne "\rIteration $i..."
	${this_dir}/../skullsup -s "${interface}" incant hellivator; sleep 1;
	${this_dir}/../skullsup -s "${interface}" color ff03f2; sleep 1;
	${this_dir}/../skullsup -s "${interface}" incant vortex; sleep 1;
	${this_dir}/../skullsup -s "${interface}" color 205020; sleep 1;
	${this_dir}/../skullsup -s "${interface}" incant pulse 306020; sleep 1;
	${this_dir}/../skullsup -s "${interface}" color 0125f2; sleep 1;
	i=$(expr $i + 1)
done
