#!/bin/sh

. common.sh

cd "$DUT_DIR"
vagrant up
cd "$OLD_PWD"
