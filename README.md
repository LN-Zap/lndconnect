# zapconnect 🌩

Generate QRCode to connect iOS app to remote LND.

## Installing zapconnect

```
go get -d github.com/LN-Zap/zapconnect
cd $GOPATH/src/github.com/LN-Zap/zapconnect
go get ./...
go install -v ./...
```

## Starting zapconnect

```
zapconnect
```

## Application Options

```
-i, --localip            Use local ip instead of public ip.
-j, --json               Display json instead of a QRCode.
    --lnddir=            The base directory that contains lnd's data, logs, configuration
                         file, etc.
    --configfile=        Path to configuration file
-b, --datadir=           The directory to store lnd's data within
    --tlscertpath=       Path to write the TLS certificate for lnd's RPC and REST services
    --adminmacaroonpath= Path to write the admin macaroon for lnd's RPC and REST services 
                         if it doesn't exist
```
