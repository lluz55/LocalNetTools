# LocalNetTools

LocalNetTools was made just to play around with GO.I wanted my PC to be able to shut down from any other device. So, I made this.

## Run

Download from release and run it.

## Build

First install *go-bindata*
```bash
go get -u github.com/shuLhan/go-bindata
```

Clone this repo and use the follow commands:
```bash
go-bindata webui/dist/...
```
```bash
go build -ldflags="-s -w -H windowsgui"
```
## Usage

```bash
./LocalNetTools.exe
```
You can pass **-s** flag to start it with Windows:

```bash
./LocalNetTools.exe -s
```
To see the gui in your default web browser click in the system tray icon and then click _Open Ui_ menu

### Important
Make sure windows firewall allow LocalNetTools to run in local network


## License
[MIT](https://github.com/lluz55/LocalNetTools/blob/master/License)
