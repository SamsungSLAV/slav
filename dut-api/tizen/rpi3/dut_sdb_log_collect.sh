#!/usr/bin/expect -f

spawn -open [open /dev/ttyS1 w+]
send "\r"
send "\r"
expect "localhost login:"
send "root\r"
expect "Password:"
send "tizen\r"

set timeout 20

expect "\# "
sleep 3
send "dlogutil SDBD_TRACE_SDB SDBD_TRACE_TRANSPORT -d | cat\r"
expect "\# "
sleep 1
send "exit\r"
close
