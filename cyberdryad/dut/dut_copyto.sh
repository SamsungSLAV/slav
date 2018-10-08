#!/bin/sh

. config.sh

SSH_CONFIG_TMP="$(mktemp)"

cleanup() {
    rm "$SSH_CONFIG_TMP"
}

get_ssh_config "$SSH_CONFIG_TMP"
scp -F "$SSH_CONFIG_TMP" "$1" "${HOSTNAME}:${2}"
cleanup
