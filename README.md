# maint-http

Fast-superlight http "backup" backend for haproxy  
This is designed to serve static http pages for maintenance purposes.

## Usage

```shell
Usage: ./maint-http [-d <directory_to_serve>] [-p <tcp_listen_port>]
```

## Note

This application binds to localhost only!

## Performance test

Use [bombardier](github.com/codesenberg/bombardier)

```shell
./bombardier -c 125 -n 2690701 http://localhost:3000
```

## Compile

```shell
go build maint-http.go
```

## HTML static pages and resource files

subfolder ```html``` contains resource files:

* index.html --> maintenance page
* ```res``` subfolder
    * *.woff2 --> fonts
    * *.jpg --> images

CSS is embedded in the html file

Of course you can totally customize your maintenance page. This is just a clean example.

## Systemd integration

Use the included systemd unit

## HAProxy configuration example

Here an example of haproxy configuration

```ini
backend soa12qa_be
    balance     roundrobin
    cookie HAPROXYSRV insert nocache
    #option httpchk
    server  soa12-qa-wt01 10.211.9.12:8777 cookie soa12-qa-wt01 check
    server  soa12-qa-wt02 10.211.9.13:8777 cookie soa12-qa-wt02 check
    server  maintenance 127.0.0.1:3000 check backup
```
