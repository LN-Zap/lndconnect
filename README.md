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

With option `-j` you can also display a json string instead of a QRCode.
