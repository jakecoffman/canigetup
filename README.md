# canigetup

Tell your kids to stay in bed when an LED is on.

## Systemd

```
scp canigetup.service pi@192.168.1.14:canigetup.service
ssh pi@192.168.1.14
sudo mv canigetup.service /lib/systemd/system
sudo chmod 644 /lib/systemd/system/canigetup.service
sudo systemctl daemon-reload
sudo systemctl enable canigetup.service
```

Then start the service or reboot to test it. 

## Build & Deploy

```
GOARCH=arm go build
npm run build
scp canigetup pi@192.168.1.14:canigetup2
scp -r dist pi@192.168.1.14:web2
ssh pi@192.168.1.14
sudo systemctl stop canigetup.service
mv canigetup2 canigetup
rm -rf web
mv web2 web
sudo systemctl start canigetup.service
exit
```
