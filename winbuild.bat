set GOOS=linux
set GOARCH=arm
go build
scp canigetup pi@192.168.1.14
