#!/bin/sh
# Setup LXC container to serve as DUT supervisor
#
# See README.dut-s for help.

set -ex

VM_NAME=${DUT_VM_NAME:-dut-s}
BORUTA_ADDRESS=${BORUTA_ADDRESS:-127.0.0.1}
BORUTA_PORT=${BORUTA_PORT:-7175}

lxc_install() {
    lxc-destroy -f -n "$VM_NAME" || true
    lxc-create -n "$VM_NAME" -t ubuntu -- -r xenial --packages=virtinst,libvirt-bin,libvirt-daemon,qemu-kvm,qemu-utils,openssh-server,wget,expect,git,libsystemd-dev,sshfs
    # print big scarry warning if network is not setup correctly
    if ! grep -q 'lxc.network.type.*veth' "/var/lib/lxc/$VM_NAME/config"; then
        echo "@@@ warning: this script assumes lxc uses veth<->lxcbr0 networking setup BY DEFAULT @@@"
        echo "if you are using different distro than ubuntu please ensure lxc.network.* is set up correctly"
    fi
}

tools_install()
{
    git archive HEAD | lxc-attach -n "$VM_NAME" -- tar xvfC - "$1/bin"
}

tools_setup()
{
    cat <<EOF | lxc-attach -n "$VM_NAME" dd of=/var/tmp/dut-s-setup.env
DESTDIR=${DESTDIR}
BORUTA_ADDRESS=${BORUTA_ADDRESS}
BORUTA_PORT=${BORUTA_PORT}
EOF
    lxc-attach -n "$VM_NAME" -- "dut-s-reconfigure.sh"
}

DESTDIR="/usr/local"

lxc_install

lxc-start -n "$VM_NAME"
# Hack: wait a while for container to boot and bring up networking - needed for tools setup
sleep 5
tools_install "$DESTDIR"
tools_setup

lxc-stop -r -n "$VM_NAME" || true

echo "DUT-supervisor successfully created as lxc $VM_NAME container"
