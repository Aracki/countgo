# countgo
Unique site visits in Go

## Usage:

Make config.yml in /etc/countgo/ with mongodb properties:
```
 host:
 database: 
 username:
 password: 
```

Download all libraries into your _$GOPATH_:

`go get ./...`

## Run:

To run application: `./go/bin/main 2>&1 >> logfile&`
