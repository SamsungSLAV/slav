#!/bin/sh

sed -i'' -e 's/192.168.122/10.0.0/g' /etc/libvirt/qemu/networks/default.xml
systemctl restart libvirtd.service
