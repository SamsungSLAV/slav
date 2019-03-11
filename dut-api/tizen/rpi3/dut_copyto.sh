#!/bin/sh

STDOUT_FILE=/tmp/weles_cft_stdout.txt
STDERR_FILE=/tmp/weles_cft_stderr.txt

rm $STDOUT_FILE $STDERR_FILE

(((sdb push $1 $2 | tee $STDOUT_FILE) 3>&1 1>&2 2>&3 | tee $STDERR_FILE) 3>&1 1>&2 2>&3) 

if [ -s $STDERR_FILE ]; then
	exit 1
else
	exit 0
fi
