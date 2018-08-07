# countgo
Unique site visits in Go

## Usage:

Generate mongo_config.yml in root directory with following mongo properties:
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

UPDATE (7.8.2018): run `go run cmd/aracki/main.go 2>&1 >> logfile&` from root of project
