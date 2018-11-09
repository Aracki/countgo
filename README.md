# countgo [![Go Report Card](https://goreportcard.com/badge/github.com/aracki/countgo)](https://goreportcard.com/report/github.com/aracki/countgo)
Hobby project.
Unique site visits.
Playing with Youtube playlists/videos.

## Prepare app:

Create configuration files in the root directory according to approprate template files.

1. `mongo_config.yml`
2. `config.yml`

Download all libraries into your vendor folder:

1. `dep init`
2. `dep ensure`

## Run app:

1. Run mongodb: `mongod --dbpath data/db/`
2. Run application: `go run cmd/aracki/main.go >> logfile&` from the root of the project. 

Use flag `-m=false` to run without mongodb.
Use flag `-y=false` to run without gotube.
Ampersand will start the process in the background.
Pipe to `tee logfile` to split logs into stdout and file.
