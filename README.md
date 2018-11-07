# countgo
Unique site visits in Go. Hobby project.

## Prepare app:

Create configuration files in the root directory according to approprate template files.

1. `mongo_config.yml`
2. `config.yml`

Download all libraries into your vendor folder:

1. `dep init`
2. `dep ensure`

## Run app:

1. Run mongodb: `mongod --dbpath data/db/`
2. Run application: `go run cmd/aracki/main.go 2>&1 >> logfile&` from root of project

Use flag `-m=false` to run without mongodb.
  
