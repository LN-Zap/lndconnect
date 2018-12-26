# zapconnect 🌩

Generate QRCode to connect iOS app to remote LND.

## Installing zapconnect

```
go get -d github.com/LN-Zap/zapconnect
cd $GOPATH/src/github.com/LN-Zap/zapconnect
make
```

## Starting zapconnect

```
zapconnect
```

## Application Options

```
-i, --localip               Include local ip in QRCode
-l, --localhost             Use 127.0.0.1 for ip
    --host=                 Use specific host name
-p, --port=                 Use this port (default: 10009)
-o, --image                 Output QRCode to file
    --invoice               Use invoice macaroon
    --readonly              Use readonly macaroon
    --lnddir=               The base directory that contains lnd's data, logs, configuration
                            file, etc.
    --configfile=           Path to configuration file
-b, --datadir=              The directory to find lnd's data within
    --tlscertpath=          Path to read the TLS certificate from
    --adminmacaroonpath=    Path to read the admin macaroon from
    --readonlymacaroonpath= Path to read the read-only macaroon from
    --invoicemacaroonpath=  Path to read the invoice-only macaroon from

```
