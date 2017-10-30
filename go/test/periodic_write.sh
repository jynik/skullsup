#!/bin/bash
#
# Periodically write data to a queue
################################################################################

this_dir=$(dirname $(realpath "${0}"))
config="${1}"
period="${2}"

commands=(\
	"incant hellivator" \
	"color ff0000" \
	"color 00ff00" \
	"incant vortex" \
	"color 0f00f0" \
	"incant pulse" \
	"incant pulse 0040ff" \
	"incant pulse 00ffff" \
	"color 00ff00" \
)

num_commands=${#commands[@]}

if [ $# -ne 1 ]; then
	echo "Usage: $(basename ${0}) <config file>" >&2
	exit 1
fi

if [[ ! "${period}" =~ ^[0-9]*$ ]]; then
	echo "Invalid period: ${period}" >&2
	exit 1
fi

while [ 1 ]; do
	for command in "${commands[@]}"; do
		echo "[$(date)] Writing ${command}..."
		${this_dir}/../skullsup-queue-writer -p 12345 -r localhost --insecure -c ${config} ${command}

		sleep_time=$(( ( RANDOM % 10 ) + 1 ))
		echo "Sleeping for ${sleep_time} seconds..."
		sleep ${sleep_time}
	done
done
