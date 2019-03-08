# countgo 

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/0e1f8594527b414890de3e5c92d7affd)](https://app.codacy.com/app/Aracki/countgo?utm_source=github.com&utm_medium=referral&utm_content=Aracki/countgo&utm_campaign=Badge_Grade_Dashboard)
[![Go Report Card](https://goreportcard.com/badge/github.com/aracki/countgo)](https://goreportcard.com/report/github.com/aracki/countgo) 
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FAracki%2Fcountgo.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FAracki%2Fcountgo?ref=badge_shield)

Hobby project.<br>
Unique site visits with MongoDB.<br>
Playing with Youtube playlists/videos.<br>

## Prepare app

Create configuration files in the root directory according to approprate template files.

1. `mongo_config.yml`
2. `config.yml`

Download all libraries into your vendor folder:

1. `dep init`
2. `dep ensure`

## Run app

1. Run mongodb: `mongod --dbpath data/db/`
2. Run application: `go run cmd/aracki/main.go 2>> logfile&` from the root of the project. 

Use flag `-m=false` to run without mongodb.<br>
Use flag `-y=false` to run without gotube.<br>
Ampersand `&` will start the process in the background.<br>
Pipe to `tee logfile` to split logs into stdout and file.<br>

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FAracki%2Fcountgo.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FAracki%2Fcountgo?ref=badge_large)
