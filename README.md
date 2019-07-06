# Can I get up?

Tell your kids to stay in bed when an LED is on.

## Systemd

```
scp canigetup.service pi@192.168.1.15:canigetup.service
ssh pi@192.168.1.15
sudo mv canigetup.service /lib/systemd/system
sudo chmod 644 /lib/systemd/system/canigetup.service
sudo systemctl daemon-reload
sudo systemctl enable canigetup.service
```

Then start the service or reboot to test it. 

## Build & Deploy

API:

```
GOARCH=arm go build
scp canigetup pi@192.168.1.15:canigetup2
ssh pi@192.168.1.15
sudo systemctl stop canigetup.service
mv canigetup2 canigetup
sudo systemctl start canigetup.service
exit
```

UI:

```
npm run build
scp -r dist pi@192.168.1.15:web2
ssh pi@192.168.1.15
rm -rf web
mv web2 web
exit
```
