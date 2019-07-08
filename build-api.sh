#!/bin/bash
GOOS=linux GOARCH=arm go build
scp canigetup pi@192.168.1.15:canigetup2
ssh pi@192.168.1.15 <<eof
sudo systemctl stop canigetup.service
mv canigetup2 canigetup
sudo systemctl start canigetup.service
exit
eof
