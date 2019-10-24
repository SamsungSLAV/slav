#!/bin/sh

INVENTORY='inventory.yml'

create_inv_ifnonex() {
    cd ansible/
    if ! [ -e "$INVENTORY" ]; then
        ln --symbolic "${INVENTORY}.sample" "$INVENTORY"
    fi
    cd "$OLDPWD"
}

run_vagrant_env() {
    cd vagrant/
    vagrant up
    cd "$OLDPWD"
}

main() {
    create_inv_ifnonex
    run_vagrant_env
}

main
