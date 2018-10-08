#!/bin/sh

. config.sh

SSH_CONFIG_TMP="$(mktemp)"

cleanup() {
    rm "$SSH_CONFIG_TMP"
}

get_ssh_config "$SSH_CONFIG_TMP"
ssh -F "$SSH_CONFIG_TMP" "$HOSTNAME" "$@"
cleanup
