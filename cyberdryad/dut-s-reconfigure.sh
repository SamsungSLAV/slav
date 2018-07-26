#!/bin/sh

set -ex

export GO=/usr/lib/go-1.10/bin/go

ARCH="$(uname -m)"
UNITDIR="/etc/systemd/system"
UNIT2ENABLE=""

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
        -e "/\[caps\]/ a device_type = \"qemu\"\narchitecture = \"$ARCH\"\nUUID = \"$UUID\"" \
        -e 's!groups =.*!groups = [ "kvm", "libvirtd", "stm" ]!' \
        /etc/boruta/dryad.conf
}

stm_setup()
{
    muxpi_path="$1"
    install -m0644 "${muxpi_path}/sw/nanopi/stm.service" "$UNITDIR/stm.service"
    sed -i -e 's/-serve/-serve -dummy/' "$UNITDIR/stm.service"
    install -m0644 "${muxpi_path}/sw/nanopi/stm.socket" "$UNITDIR/stm.socket"
    install -m0644 "${muxpi_path}/sw/nanopi/stm-user.socket" "$UNITDIR/stm-user.socket"
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
    HOST=git.tizen.org
    D=/src/slav
    SRC="$D/src/$HOST"

    mkdir -p "$D/pkg" "$SRC/tools"

    for i in tools/boruta tools/weles tools/muxpi; do
        # assume previous invocation might have failed - better to fetch again
        test -d "$SRC/$i" && rm -rf "$SRC/$i"
        git clone "git://$HOST/$i" "$SRC/$i"
    done

    export GOPATH="$D"
    (cd "$D" && $GO get ./...)
    install -m0755 -t "$DESTDIR/bin" "$GOPATH/bin/stm" "$GOPATH/bin/dryad"

    dryad_setup
    dryad_unit_install
    stm_setup "$SRC/tools/muxpi"
    systemd_unit_enable
}

# This script assumes DESTDIR, UUID, BORUTA_ADDRESS & BORUTA_PORT
# is defined in the environment.  In normal setup it's passed from
# dut-s-lxc-setup via following file:
. /var/tmp/dut-s-setup.env

apt_setup
slav_setup
