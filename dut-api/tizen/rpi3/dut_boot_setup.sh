#!/usr/bin/expect -f

set DUT_IP [lindex $argv 0]

set timeout 60

#spawn screen -S expect_screen /dev/ttyS1 115200
spawn -open [open /dev/ttyS1 w+]
send "\r"
send "\r"
expect "localhost login:"
send "root\r"
expect "Password:"
send "tizen\r"

set timeout 1

expect "\# "
sleep 3
send "ifconfig eth0 $DUT_IP netmask 255.255.255.0\r"
expect "\# "
sleep 1
send "exit\r"
close
#exec screen -S expect_screen -X quit
