#!/bin/sh

#DUT_DIR='/home/vagrant/dut'
DUT_DIR='.'
HOSTNAME='default'

get_ssh_config() {
    local ssh_config="$1"
    cd "$DUT_DIR"
    vagrant ssh-config > "$ssh_config"
    cd "$OLD_PWD"
}
