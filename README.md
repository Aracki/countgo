# countgo 
[![Go Report Card](https://goreportcard.com/badge/github.com/aracki/countgo)](https://goreportcard.com/report/github.com/aracki/countgo) 
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FAracki%2Fcountgo.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FAracki%2Fcountgo?ref=badge_shield)

Hobby project.<br>

It contains:
* Web server serving [static](https://github.com/Aracki/countgo/tree/master/static) website with TLS. (previously - [Aracki/aracki.me](https://github.com/Aracki/aracki.me))
* Unique site visits counter with MongoDB.
* Playing with Youtube playlists/videos.

## Prepare app

Create configuration files in the root directory according to approprate template files.

1. `mongo_config.yml`
2. `config.yml`

Download all libraries into your vendor folder:

1. `dep init`
2. `dep ensure`

## Run app

1. Run mongodb: `mongod --dbpath data/db/`
2. Run application: `go run cmd/aracki/main.go -y=false >> logfile 2>&1 &` from the root of the project. 

Use flag `-m=false` to run without mongodb.<br>
Use flag `-y=false` to run without gotube.<br>
Ampersand `&` will start the process in the background.<br>
Pipe to `tee logfile` to split logs into stdout and file.<br>

## Set up TLS

* Local development:
    * Install [mkcert](https://github.com/FiloSottile/mkcert) tool
    * `mkcert -install`
    * `mkcert localhost 127.0.0.1`
* Production:
    * Install [certbot](https://github.com/certbot/certbot) for automatic obtaining/renewing Let's Encrypt certificates
    * `./certbot-auto certonly`
    

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FAracki%2Fcountgo.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FAracki%2Fcountgo?ref=badge_large)
