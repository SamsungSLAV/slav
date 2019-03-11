#!/bin/bash

stty -F /dev/ttyS1 115200 cs8 -cstopb -parenb

DUT_IP="192.168.69.10" # TODO: Obtain this address from somwhere in  he future. Maybe env?

sdb disconnect $DUT_IP > /dev/null

/usr/local/bin/stm -dut
/usr/local/bin/stm -dut # FIXME: First attempt to call stm fails. Probably fix in stm firmware and/or stm executable (golang side) is needed.
sleep 2
/usr/local/bin/stm -tick

/usr/local/bin/dut_boot_setup.sh $DUT_IP > /tmp/dut_but_setup.log
sleep 10

sdb connect $DUT_IP
CRES=$(sdb connect $DUT_IP)

/usr/local/bin/dut_sdb_log_collect.sh > "/tmp/dut_sdb_log_collect_$(date +%s).log" 2>&1

if [[ $CRES == "error:"* || $CRES == *"fail"* ]]; then
	(>&2 echo -e $CRES)
	exit 1
else
	echo $CRES
	exit 0
fi
