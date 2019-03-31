# dirhttps

Spinning up an HTTPS server from current directory - the TLS equivalent of `python -m http.server`.

## Install

Use [prebuilt binaries](https://github.com/maetthu/dirhttps/releases)

or build from source

```
go get -u github.com/maetthu/dirhttps
$(go env GOPATH)/bin/dirhttps
```

## Setup

_dirhttps_ needs a certificate and corresponding key to operate. Easiest option is to use the excellent [mkcert](https://github.com/FiloSottile/mkcert) tool which creates locally trusted, self-signed development certificates.

* Install and setup mkcert CA
* Create dirhttps certificate

``` 
$ mkdir ~/.config/dirhttps
$ cd ~/.config/dirhttps
$ mkcert -cert-file cert.pem -key-file key.pem localhost 127.0.0.1 more.hostnames.or.ips.if.needed.example.org
```

* Done

## Usage

```
$ cd /tmp
$ dirhttps  
2019/03/31 13:40:26 Listening for HTTPS connections on :8443
2019/03/31 13:40:26 Serving from directory /tmp
```

### Different listening address/port

```
$ dirhttps -l :1234
$ dirhttps -l 127.0.0.2:8443
```

