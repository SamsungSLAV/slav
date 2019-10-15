#!/bin/sh

set -ex

export GO=/usr/lib/go-1.10/bin/go

RELEASE='v0.1.0'

ARCH="$(uname -m)"
UNITDIR="/etc/systemd/system"
UNIT2ENABLE=""

storage_setup()
{
    mkdir -m 1777 /var/dut
}

net_setup()
{
    cat > "$UNITDIR/net-tun.service" <<EOF
[Unit]
Description=/dev/net/tun for libvirt/qemu
Before=libvirtd.service

[Service]
Type=oneshot
ExecStartPre=-/bin/mkdir -m 0755 /dev/net
ExecStart=/bin/mknod -m 0666 /dev/net/tun c 10 200

[Install]
WantedBy=default.target
EOF
    UNITS2ENABLE="$UNITS2ENABLE net-tun.service"
    systemctl start net-tun.service

    virsh net-define --file /dev/stdin <<EOF
<network>
  <name>dut-network</name>
  <bridge name='dut-br0' stp='on' delay='0'/>
  <domain name='dut-network'/>
  <ip address='192.168.100.1' netmask='255.255.255.0'>
    <dhcp>
      <range start='192.168.100.100' end='192.168.100.100'/>
    </dhcp>
  </ip>
</network>
EOF
   virsh net-autostart dut-network
   virsh net-start dut-network
}

apt_setup()
{
    echo 'deb http://archive.ubuntu.com/ubuntu xenial-backports main restricted universe multiverse' > /etc/apt/sources.list.d/dut-s-xenial-backports.list
    apt-get update
    apt-get install -y golang-1.10-go
}

dryad_setup()
{
    mkdir -p /etc/boruta
    # dryad has to fail in first run AND create configuration template
    # we run it here for the template
    dryad || true
    sed -i \
        -e 's,boruta_address =.*,boruta_address = "'${BORUTA_ADDRESS}:${BORUTA_PORT}'",' \
        -e "/\[caps\]/ a device_type = \"qemu\"\nDeviceType = \"qemu\"\narchitecture = \"$ARCH\"\nUUID = \"$UUID\"" \
        -e 's!groups =.*!groups = [ "kvm", "libvirtd", "stm" ]!' \
        /etc/boruta/dryad.conf
}

stm_setup()
{
    muxpi_path="$1"
    install -m0644 "${muxpi_path}/sw/nanopi/stm/systemd/stm.service" "$UNITDIR/stm.service"
    sed -i -e 's/stm -serve/stm.real -serve -dummy/' "$UNITDIR/stm.service"
    install -m0644 "${muxpi_path}/sw/nanopi/stm/systemd/stm.socket" "$UNITDIR/stm.socket"
    install -m0644 "${muxpi_path}/sw/nanopi/stm/systemd/stm-user.socket" "$UNITDIR/stm-user.socket"
    sed -i -e "s,/usr/bin,$DESTDIR/bin,g" -e "s,WantedBy=sockets.target,WantedBy=basic.target," "$UNITDIR/stm.service" "$UNITDIR/stm.socket" "$UNITDIR/stm-user.socket"
    UNITS2ENABLE="$UNITS2ENABLE stm.socket stm-user.socket"

    getent group stm || addgroup --system stm
}

dryad_unit_install()
{
    cat > "$UNITDIR/dryad.service" <<EOF
[Unit]
Description=SLAV Dryad - Boruta Server Agent
Requires=stm.socket

[Service]
Type=simple
ExecStartPre=-/bin/mknod -m 0666 /dev/fuse c 10 229
ExecStart=$DESTDIR/bin/dryad
Restart=always

[Install]
WantedBy=default.target
EOF
    UNITS2ENABLE="$UNITS2ENABLE dryad.service"
}

systemd_unit_enable()
{
    systemctl daemon-reload
    systemctl enable $UNITS2ENABLE
}

slav_setup()
{
    HOST=github.com
    D=/src/slav
    SRC="$D/src/$HOST"

    mkdir -p "$D/pkg" "$SRC/SamsungSLAV"

    for i in SamsungSLAV/boruta SamsungSLAV/weles SamsungSLAV/muxpi; do
        # assume previous invocation might have failed - better to fetch again
        test -d "$SRC/$i" && rm -rf "$SRC/$i"
        git clone --branch "$RELEASE" "https://$HOST/$i.git" "$SRC/$i"
    done

    export GOPATH="$D"
    (cd "$D" && $GO get ./...)
    install -m0755 -t "$DESTDIR/bin" "$SRC/SamsungSLAV/muxpi/sw/nanopi/stm/stm" "$GOPATH/bin/dryad"
    install -m0755 "$GOPATH/bin/stm" "$DESTDIR/bin/stm.real"
    sed -i -e "s,/usr/bin/stm,$DESTDIR/bin/stm.real," "$DESTDIR/bin/stm"

    dryad_setup
    dryad_unit_install
    stm_setup "$SRC/SamsungSLAV/muxpi"
    systemd_unit_enable
}

# This script assumes DESTDIR, UUID, BORUTA_ADDRESS & BORUTA_PORT
# is defined in the environment.  In normal setup it's passed from
# dut-s-lxc-setup via following file:
. /var/tmp/dut-s-setup.env

net_setup
apt_setup
storage_setup
slav_setup
